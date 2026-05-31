package memory_dto

type GetMemoryStatsInput struct {
}

type GetMemoryStatsOutput struct {
	Stats MemoryStatsItem `json:"stats"`
}
