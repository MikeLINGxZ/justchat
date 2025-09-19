package main

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/schema"
	agents "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/example/internal"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
)

func main() {

	is, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	openAIAPIKey := os.Getenv("ALICCLOUD_API_KEY")
	openAIBaseURL := "https://dashscope.aliyuncs.com/compatible-mode/v1"
	openAIModelName := "qwen-flash"

	ctx := context.Background()
	h, err := internal.NewHost(ctx, openAIBaseURL, openAIAPIKey, openAIModelName)
	if err != nil {
		panic(err)
	}

	memoryAgent, err := agents.NewMemoryAgent(ctx, openAIBaseURL, openAIAPIKey, openAIModelName, is)
	if err != nil {
		panic(err)
	}

	hostMA, err := host.NewMultiAgent(ctx, &host.MultiAgentConfig{
		Host: *h,
		Specialists: []*host.Specialist{
			memoryAgent,
		},
	})
	if err != nil {
		panic(err)
	}

	cb := &logCallback{}

	for { // 多轮对话，除非用户输入了 "exit"，否则一直循环
		println("\n\nYou: ") // 提示轮到用户输入了

		var message string
		scanner := bufio.NewScanner(os.Stdin) // 获取用户在命令行的输入
		for scanner.Scan() {
			message += scanner.Text()
			break
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		if message == "exit" {
			return
		}

		msg := &schema.Message{
			Role:    schema.User,
			Content: message,
		}

		out, err := hostMA.Stream(ctx, []*schema.Message{msg}, host.WithAgentCallbacks(cb))
		if err != nil {
			panic(err)
		}

		println("\nAnswer:")

		for {
			msg, err := out.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
			}

			print(msg.Content)
		}

		out.Close()
	}
}

type logCallback struct{}

func (l *logCallback) OnHandOff(ctx context.Context, info *host.HandOffInfo) context.Context {
	println("\nHandOff to", info.ToAgentName, "with argument", info.Argument)
	return ctx
}
