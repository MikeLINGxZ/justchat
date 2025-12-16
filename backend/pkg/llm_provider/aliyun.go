package llm_provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type Aliyun struct {
	providerModel wrapper_models.ProviderModel
}

func NewAliyun(providerModel wrapper_models.ProviderModel) IProvider {
	return &Aliyun{
		providerModel: providerModel,
	}
}

func (a *Aliyun) Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		BaseURL: a.providerModel.BaseUrl,
		Model:   a.providerModel.Model,
		APIKey:  a.providerModel.ApiKey,
	})
	if err != nil {
		return nil, err
	}

	var messagesPoint []*schema.Message
	for _, item := range messages {
		messagesPoint = append(messagesPoint, &item)
	}

	// 调用LLM服务
	streamResult, err := chatModel.Stream(ctx, messagesPoint)
	if err != nil {
		return nil, err
	}

	return streamResult, nil
}

func (a *Aliyun) BuildUserMessage(ctx context.Context, message view_models.MessagePkg) (*schema.Message, error) {
	var paths []string
	path2Url := make(map[string]string)
	for _, file := range message.Files {
		paths = append(paths, file.FilePath)
	}
	for _, path := range paths {
		uploadPolicy, err := a.getUploadPolicy()
		if err != nil {
			return nil, err
		}
		fileUrl, err := a.uploadFileToOSS(uploadPolicy, path)
		if err != nil {
			return nil, err
		}
		path2Url[path] = fileUrl
	}

	var userInputMultiContent []schema.MessageInputPart
	if message.Content != "" {
		userInputMultiContent = append(userInputMultiContent, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeText,
			Text: message.Content,
		})
	}

	for _, item := range message.Files {
		var text string
		var img *schema.MessageInputImage
		var audio *schema.MessageInputAudio
		var video *schema.MessageInputVideo
		var file *schema.MessageInputFile
		url := path2Url[item.FilePath]

		messagePartCommon := schema.MessagePartCommon{
			URL:        &url,
			Base64Data: nil,
			MIMEType:   item.MineType,
			Extra: map[string]interface{}{
				"path":      url,
				"mime_type": item.MineType,
			},
		}
		switch item.ChatMessagePartType {
		case schema.ChatMessagePartTypeText:
			continue
		case schema.ChatMessagePartTypeImageURL:
			img = &schema.MessageInputImage{
				MessagePartCommon: messagePartCommon,
				Detail:            schema.ImageURLDetailHigh,
			}
		case schema.ChatMessagePartTypeAudioURL:
			audio = &schema.MessageInputAudio{
				MessagePartCommon: messagePartCommon,
			}
		case schema.ChatMessagePartTypeVideoURL:
			video = &schema.MessageInputVideo{
				MessagePartCommon: messagePartCommon,
			}
		case schema.ChatMessagePartTypeFileURL:
			file = &schema.MessageInputFile{
				MessagePartCommon: messagePartCommon,
			}
		}
		if img == nil && audio == nil && video == nil && file == nil {
			continue
		}
		userInputMultiContent = append(userInputMultiContent, schema.MessageInputPart{
			Type:  item.ChatMessagePartType,
			Text:  text,
			Image: img,
			Audio: audio,
			Video: video,
			File:  file,
		})
	}

	return &schema.Message{
		Role:                  schema.User,
		Content:               "",
		UserInputMultiContent: userInputMultiContent,
	}, nil
}

// AliyunUploadPolicyResponse 上传凭证响应结构
type AliyunUploadPolicyResponse struct {
	Data AliyunUploadPolicyData `json:"data"`
}

// AliyunUploadPolicyData 上传凭证数据
type AliyunUploadPolicyData struct {
	UploadDir           string `json:"upload_dir"`
	OssAccessKeyId      string `json:"oss_access_key_id"`
	Signature           string `json:"signature"`
	Policy              string `json:"policy"`
	XOssObjectAcl       string `json:"x_oss_object_acl"`
	XOssForbidOverwrite string `json:"x_oss_forbid_overwrite"`
	UploadHost          string `json:"upload_host"`
}

// getUploadPolicy 获取文件上传凭证
func (a *Aliyun) getUploadPolicy() (*AliyunUploadPolicyData, error) {
	url := *a.providerModel.FileUploadBaseUrl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.providerModel.ApiKey))
	req.Header.Set("Content-Type", "application/json")

	// 添加查询参数
	q := req.URL.Query()
	q.Add("action", "getPolicy")
	q.Add("model", a.providerModel.Model)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
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
		return nil, fmt.Errorf("failed to get upload policy: status %d, body: %s", resp.StatusCode, string(body))
	}

	var policyResp AliyunUploadPolicyResponse
	if err := json.Unmarshal(body, &policyResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &policyResp.Data, nil
}

// uploadFileToOSS 将文件上传到临时存储OSS
func (a *Aliyun) uploadFileToOSS(policyData *AliyunUploadPolicyData, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	key := fmt.Sprintf("%s/%s", policyData.UploadDir, fileName)

	// 创建 multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加表单字段
	fields := map[string]string{
		"OSSAccessKeyId":         policyData.OssAccessKeyId,
		"Signature":              policyData.Signature,
		"policy":                 policyData.Policy,
		"x-oss-object-acl":       policyData.XOssObjectAcl,
		"x-oss-forbid-overwrite": policyData.XOssForbidOverwrite,
		"key":                    key,
		"success_action_status":  "200",
	}

	for fieldName, fieldValue := range fields {
		if err := writer.WriteField(fieldName, fieldValue); err != nil {
			return "", fmt.Errorf("failed to write field %s: %w", fieldName, err)
		}
	}

	// 添加文件
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", policyData.UploadHost, &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload file: status %d, body: %s", resp.StatusCode, string(body))
	}

	return fmt.Sprintf("oss://%s", key), nil
}
