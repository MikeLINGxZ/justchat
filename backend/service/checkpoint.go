package service

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

type CheckpointStage string

const (
	CheckpointStageClassify   CheckpointStage = "classify"
	CheckpointStagePlan       CheckpointStage = "plan"
	CheckpointStageBatchDone  CheckpointStage = "batch_done"
	CheckpointStageSynthesize CheckpointStage = "synthesize"
	CheckpointStageReview     CheckpointStage = "review"
	CheckpointStageFinalize   CheckpointStage = "finalize"
)

type CheckpointData struct {
	Stage          CheckpointStage               `json:"stage"`
	Plan           *workflowPlan                 `json:"plan,omitempty"`
	Results        map[string]workflowTaskResult `json:"results,omitempty"`
	Draft          string                        `json:"draft,omitempty"`
	RetryCount     int                           `json:"retry_count,omitempty"`
	ReviewFeedback string                        `json:"review_feedback,omitempty"`
	Timestamp      time.Time                     `json:"timestamp"`
}

type CheckpointManager struct {
	runner *completionRunner
}

func NewCheckpointManager(runner *completionRunner) *CheckpointManager {
	return &CheckpointManager{runner: runner}
}

func (m *CheckpointManager) Save(ctx context.Context, stage CheckpointStage, data CheckpointData) error {
	data.Stage = stage
	data.Timestamp = time.Now()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	m.runner.mu.Lock()
	if m.runner.assistantMessage.AssistantMessageExtra == nil {
		m.runner.assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
	}
	m.runner.assistantMessage.AssistantMessageExtra.CurrentStage = string(stage)
	m.runner.mu.Unlock()

	_ = jsonData
	return nil
}

func (m *CheckpointManager) Load(ctx context.Context) (*CheckpointData, error) {
	_ = ctx
	return nil, nil
}

func (m *CheckpointManager) SetupCheckpointRecovery(
	ctx context.Context,
	runner *completionRunner,
	handoff *workflowHandoff,
) (recoveredPlan *workflowPlan, recoveredResults map[string]workflowTaskResult, shouldSkipTo string) {
	cp, err := m.Load(ctx)
	if err != nil || cp == nil {
		return nil, nil, ""
	}

	logger.Info("checkpoint recovery: found checkpoint at stage ", string(cp.Stage))

	switch cp.Stage {
	case CheckpointStagePlan:
		return nil, nil, "plan"
	case CheckpointStageBatchDone:
		if cp.Plan != nil {
			if cp.Results == nil {
				cp.Results = make(map[string]workflowTaskResult)
			}
			return cp.Plan, cp.Results, "synthesize"
		}
		return nil, nil, ""
	case CheckpointStageSynthesize:
		if cp.Plan != nil {
			return cp.Plan, cp.Results, "synthesize"
		}
		return nil, nil, ""
	case CheckpointStageReview:
		if cp.Plan != nil {
			return cp.Plan, cp.Results, "review"
		}
		return nil, nil, ""
	case CheckpointStageFinalize:
		return cp.Plan, cp.Results, "finalize"
	default:
		return nil, nil, ""
	}
}
