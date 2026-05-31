package memory_dto

type ListMemoriesInput struct {
	Query            string `json:"query"`
	Type             string `json:"type"`
	Target           string `json:"target"`
	IncludeForgotten bool   `json:"include_forgotten"`
	Offset           int    `json:"offset"`
	Limit            int    `json:"limit"`
}

type ListMemoriesOutput struct {
	Items []MemoryItem `json:"items"`
	Total int64        `json:"total"`
}
