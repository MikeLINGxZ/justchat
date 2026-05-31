package agent_dto

// AttachmentInput 是前端发送消息时携带的附件元数据。
// Name/Mime 可空，后端通过 path 推断后写入存储与 model.Message。
type AttachmentInput struct {
	Path string `json:"path"`
	Name string `json:"name,omitempty"`
	Mime string `json:"mime,omitempty"`
}
