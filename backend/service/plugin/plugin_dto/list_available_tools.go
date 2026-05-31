package plugin_dto

type ToolItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type ListAvailableToolsInput struct{}

type ListAvailableToolsOutput struct {
	Tools []ToolItem `json:"tools"`
}
