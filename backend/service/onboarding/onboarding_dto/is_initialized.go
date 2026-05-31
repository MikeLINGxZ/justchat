package onboarding_dto

// IsInitializedInput is the empty input for IsInitialized.
type IsInitializedInput struct{}

// IsInitializedOutput carries whether the application has finished onboarding.
type IsInitializedOutput struct {
	Initialized bool `json:"initialized"`
}
