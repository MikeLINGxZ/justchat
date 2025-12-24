package llm_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type IProvider interface {
	Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error)
	BuildUserMessage(ctx context.Context, message view_models.MessagePkg) (*schema.Message, error)
}

func NewLlmProvider(providerModel wrapper_models.ProviderModel) IProvider {
	var iProvider IProvider
	switch providerModel.ProviderType {
	case data_models.ProviderTypeDeepseek:
		iProvider = NewDeepseek(providerModel)
	case data_models.ProviderTypeAliyuns:
		iProvider = NewAliyun(providerModel)
	case data_models.ProviderTypeOpenrouter:
		iProvider = NewOpenrouter(providerModel)
	default:
		iProvider = NewOpenai(providerModel)
	}
	return iProvider
}

// ModelData represents the structure for model information
type ModelData struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
	Object  string `json:"object"`
}

// ModelsResponse represents the response structure from the models API
type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelData `json:"data"`
}

// GetModels retrieves available models from the LLM provider
func GetModels(baseURL, apiKey string) ([]ModelData, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/models", baseURL)

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var modelsResp ModelsResponse
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return modelsResp.Data, nil
}

// GenChatTitle 生成一个聊天的标题
func GenChatTitle(provider IProvider, messages []schema.Message) (string, error) {
	genTitle := ""
	contextMessages := []schema.Message{
		{
			Role: schema.System,
			Content: `
					你是一位专业的对话摘要与标题提炼专家。请根据我提供的聊天记录，生成1个最合适的标题，要求满足以下所有条件：
					✅ 准确概括核心主题：抓住双方讨论的实质焦点（如问题、决策、情感、事件或共识），而非罗列细节；
					✅ 简洁有力：控制在8–15个汉字以内，避免标点（除必要顿号）、英文和冗余修饰；
					✅ 中性客观，不带主观判断或情绪渲染（除非聊天本身是明确的情感倾诉，此时可适度体现温度，如“深夜倾诉：关于成长的迷茫与自我接纳”）；
					✅ 适配通用场景：标题应便于归档、检索或快速理解，不依赖上下文即可读懂；
					✅ 直接输出标题，不需要其他内容；
					❌ 不要解释、不要复述对话、不要添加额外信息、不要输出任何说明文字——只输出标题本身，且仅一行。
					
					请严格遵循以上规则。现在，我的聊天记录如下：
					`,
		},
	}
	contextMessages = append(contextMessages, messages...)
	resp, err := provider.Completions(context.Background(), contextMessages)
	if err == nil {
		for {
			recv, err := resp.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", err
			}
			genTitle += recv.Content
		}
	}
	if genTitle == "" {
		return "", fmt.Errorf("failed to generate title")
	}
	return genTitle, nil
}
