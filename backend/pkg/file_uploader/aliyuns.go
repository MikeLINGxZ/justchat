package file_uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type Aliyuns struct {
	providerModel *wrapper_models.ProviderModel
}

func (a *Aliyuns) Upload(paths []string) (map[string]string, error) {
	result := make(map[string]string)

	// 获取上传凭证，上传凭证接口有限流，超出限流将导致请求失败
	policyData, err := a.getUploadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to get upload policy: %w", err)
	}

	// 上传每个文件到OSS
	for _, path := range paths {
		ossURL, err := a.uploadFileToOSS(policyData, path)
		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", path, err)
		}
		result[path] = ossURL
	}

	return result, nil
}

// UploadPolicyResponse 上传凭证响应结构
type UploadPolicyResponse struct {
	Data UploadPolicyData `json:"data"`
}

// UploadPolicyData 上传凭证数据
type UploadPolicyData struct {
	UploadDir           string `json:"upload_dir"`
	OssAccessKeyId      string `json:"oss_access_key_id"`
	Signature           string `json:"signature"`
	Policy              string `json:"policy"`
	XOssObjectAcl       string `json:"x_oss_object_acl"`
	XOssForbidOverwrite string `json:"x_oss_forbid_overwrite"`
	UploadHost          string `json:"upload_host"`
}

// getUploadPolicy 获取文件上传凭证
func (a *Aliyuns) getUploadPolicy() (*UploadPolicyData, error) {
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

	var policyResp UploadPolicyResponse
	if err := json.Unmarshal(body, &policyResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &policyResp.Data, nil
}

// uploadFileToOSS 将文件上传到临时存储OSS
func (a *Aliyuns) uploadFileToOSS(policyData *UploadPolicyData, filePath string) (string, error) {
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
