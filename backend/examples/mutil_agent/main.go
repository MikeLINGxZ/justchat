// 使用 Eino 框架实现的 AI 聊天应用
// Agent-to-Agent 架构：Chat 本身是主 Agent，通过 AgentTool 委托给 DateTimeAgent、FruitPriceAgent

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/examples/mutil_agent/agents"
)

const (
	providerQwen   = "qwen"
	providerOllama = "ollama"
	providerOpenAI = "openai"
	defaultModel   = "llama3.2"
	ollamaBaseURL  = "http://localhost:11434"
	qwenBaseURL    = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	qwenDefault    = "qwen-plus-latest"
	exitCommand    = "/exit"
	clearCommand   = "/clear"
	helpCommand    = "/help"
)

func main() {
	ctx := context.Background()

	chatModel, provider, modelName := initChatModel(ctx)
	if chatModel == nil {
		log.Fatal("无法初始化 ChatModel，请检查配置")
	}

	var runner *adk.Runner
	if tcm, ok := chatModel.(model.ToolCallingChatModel); ok {
		// Agent-to-Agent：Chat 为主 Agent，子 Agent 作为 AgentTool
		dateTimeAgent, err := agents.NewDateTimeAgent(ctx, tcm)
		if err != nil {
			log.Printf("日期时间 Agent 初始化失败: %v\n", err)
		}
		fruitAgent, err := agents.NewFruitPriceAgent(ctx, tcm)
		if err != nil {
			log.Printf("水果价格 Agent 初始化失败: %v\n", err)
		}

		subAgents := make([]adk.Agent, 0, 2)
		if dateTimeAgent != nil {
			subAgents = append(subAgents, dateTimeAgent)
		}
		if fruitAgent != nil {
			subAgents = append(subAgents, fruitAgent)
		}

		chatAgent, err := agents.NewChatAgent(ctx, tcm, subAgents)
		if err != nil {
			log.Fatalf("Chat Agent 初始化失败: %v\n", err)
		}

		runner = adk.NewRunner(ctx, adk.RunnerConfig{
			Agent:           chatAgent,
			EnableStreaming: true,
		})
	} else {
		log.Fatal("当前模型不支持 Tool Calling，无法使用 Agent 架构。请使用 Qwen、OpenAI 或支持 Tool Calling 的 Ollama 模型。")
	}

	fmt.Printf("\n🤖 AI Chat 已启动 (Provider: %s, Model: %s)\n", provider, modelName)
	fmt.Printf("架构: ChatAgent → [DateTimeAgent | FruitPriceAgent]\n")
	fmt.Printf("命令: %s 退出 | %s 清空 | %s 帮助\n\n", exitCommand, clearCommand, helpCommand)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("你: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\n再见！")
				return
			}
			log.Printf("读取输入失败: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		switch input {
		case exitCommand:
			fmt.Println("再见！")
			return
		case clearCommand:
			// Agent 模式每次 Query 独立，清空仅提示
			fmt.Println("（下次对话将开启新会话）")
			continue
		case helpCommand:
			fmt.Printf("命令: %s 退出 | %s 清空\n", exitCommand, clearCommand)
			fmt.Println("直接输入问题即可，ChatAgent 会自动决定是否调用子 Agent：")
			fmt.Println("  - 日期时间 → DateTimeAgent  - 水果价格 → FruitPriceAgent  - 其他 → 直接回答")
			continue
		}

		runAgent(ctx, runner, input)
	}
}

func runAgent(ctx context.Context, runner *adk.Runner, input string) {
	fmt.Print("AI: ")
	iter := runner.Query(ctx, input)
	var reply strings.Builder

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Printf("\nAgent 执行失败: %v\n", event.Err)
			return
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		mo := event.Output.MessageOutput
		if mo.Role != schema.Assistant {
			continue
		}
		if mo.IsStreaming && mo.MessageStream != nil {
			for {
				msg, err := mo.MessageStream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("\n流式读取失败: %v\n", err)
					return
				}
				if msg != nil && msg.Content != "" {
					fmt.Print(msg.Content)
					reply.WriteString(msg.Content)
				}
			}
			mo.MessageStream.Close()
		} else if mo.Message != nil && mo.Message.Content != "" {
			fmt.Print(mo.Message.Content)
			reply.WriteString(mo.Message.Content)
		}
	}
	replyStr := reply.String()
	if replyStr == "" {
		replyStr = "(无回复)"
		fmt.Print(replyStr)
	}
	fmt.Println()
	fmt.Println()
}

func initChatModel(ctx context.Context) (model.BaseChatModel, string, string) {
	aliyunKey := os.Getenv("ALIYUN_API_KEY")
	if aliyunKey != "" {
		modelName := os.Getenv("QWEN_MODEL")
		if modelName == "" {
			modelName = qwenDefault
		}
		baseURL := os.Getenv("QWEN_BASE_URL")
		if baseURL == "" {
			baseURL = qwenBaseURL
		}
		chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
			APIKey:  aliyunKey,
			BaseURL: baseURL,
			Model:   modelName,
		})
		if err != nil {
			log.Printf("Qwen 初始化失败: %v\n", err)
		} else {
			return chatModel, providerQwen, modelName
		}
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		modelName := os.Getenv("OPENAI_MODEL")
		if modelName == "" {
			modelName = "gpt-4o-mini"
		}
		chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
			Model:  modelName,
			APIKey: apiKey,
		})
		if err != nil {
			log.Printf("OpenAI 初始化失败: %v\n", err)
		} else {
			return chatModel, providerOpenAI, modelName
		}
	}

	modelName := os.Getenv("MODEL_NAME")
	if modelName == "" {
		modelName = defaultModel
	}
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = ollamaBaseURL
	}

	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
	})
	if err != nil {
		log.Printf("Ollama 初始化失败: %v\n", err)
		log.Println("请确保 Ollama 已安装并运行: https://ollama.com")
		log.Println("或设置 ALIYUN_API_KEY 使用 Qwen / OPENAI_API_KEY 使用 OpenAI")
		return nil, "", ""
	}

	return chatModel, providerOllama, modelName
}
