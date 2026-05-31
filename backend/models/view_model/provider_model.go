package view_model

type ProviderModel struct {
	Provider Provider `json:"provider"`
	Models   []Model  `json:"models"`
}
