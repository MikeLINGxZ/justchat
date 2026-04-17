package service

import (
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

func TestUpdateCustomAgentPersistsEditableFields(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	svc := NewService()
	created, err := svc.CreateCustomAgent(view_models.CustomAgentInput{
		ID:          "custom-save-test",
		Name:        "Original Name",
		Description: "Original Description",
		Prompt:      "Original Prompt",
		Tools:       []string{"time"},
		Skills:      []string{"writer"},
	})
	if err != nil {
		t.Fatalf("CreateCustomAgent() error = %v", err)
	}
	if created == nil {
		t.Fatal("CreateCustomAgent() = nil, want detail")
	}

	updated, err := svc.UpdateCustomAgent(view_models.CustomAgentInput{
		ID:          "custom-save-test",
		Name:        "Updated Name",
		Description: "Updated Description",
		Prompt:      "Updated Prompt",
		Tools:       []string{"weather"},
		Skills:      []string{"editor"},
	})
	if err != nil {
		t.Fatalf("UpdateCustomAgent() error = %v", err)
	}
	if updated == nil {
		t.Fatal("UpdateCustomAgent() = nil, want detail")
	}

	if updated.DisplayName != "Updated Name" {
		t.Fatalf("UpdateCustomAgent().DisplayName = %q, want %q", updated.DisplayName, "Updated Name")
	}
	if updated.Description != "Updated Description" {
		t.Fatalf("UpdateCustomAgent().Description = %q, want %q", updated.Description, "Updated Description")
	}
	if got := updated.Prompts[0].Content; got != "Updated Prompt" {
		t.Fatalf("UpdateCustomAgent().Prompts[0].Content = %q, want %q", got, "Updated Prompt")
	}

	reloaded, err := svc.GetAgent("custom-save-test")
	if err != nil {
		t.Fatalf("GetAgent() error = %v", err)
	}
	if reloaded == nil {
		t.Fatal("GetAgent() = nil, want detail")
	}

	if reloaded.DisplayName != "Updated Name" {
		t.Fatalf("GetAgent().DisplayName = %q, want %q", reloaded.DisplayName, "Updated Name")
	}
	if reloaded.Description != "Updated Description" {
		t.Fatalf("GetAgent().Description = %q, want %q", reloaded.Description, "Updated Description")
	}
	if got := reloaded.Prompts[0].Content; got != "Updated Prompt" {
		t.Fatalf("GetAgent().Prompts[0].Content = %q, want %q", got, "Updated Prompt")
	}
}
