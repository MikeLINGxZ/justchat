package provider

type Type string

const (
	Deepseek            Type = "deepseek"
	Aliyun              Type = "aliyun"
	Openrouter          Type = "openrouter"
	Ollama              Type = "ollama"
	OpenAiCompatibility Type = "openai_compatibility"
)

func (p Type) String() string {
	return string(p)
}
