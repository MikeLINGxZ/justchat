package memory_dto

type GetMemoryInput struct {
	ID uint `json:"id"`
}

type GetMemoryOutput struct {
	Memory MemoryItem `json:"memory"`
}
