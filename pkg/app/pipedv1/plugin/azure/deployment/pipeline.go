package deployment

import (
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"slices"
)

type stage string

const (
	AzureFuncSync     stage = "AZURE_FUNCTION_SYNC"
	AzureFuncSwap     stage = "AZURE_FUNCTION_SWAP"
	AzureFuncRollback stage = "AZURE_FUNCTION_ROLLBACK"
)

var allStages = []string{
	string(AzureFuncSync),
	string(AzureFuncSwap),
	string(AzureFuncRollback),
}

func buildQuickSync(autoRollback bool) []sdk.QuickSyncStage {
	out := make([]sdk.QuickSyncStage, 0, 2)
	out = append(out, sdk.QuickSyncStage{
		Name:               string(AzureFuncSync),
		Description:        "", //TODO: add description
		Metadata:           map[string]string{},
		AvailableOperation: sdk.ManualOperationNone,
	})
	if autoRollback {
		out = append(out, sdk.QuickSyncStage{
			Name:               string(AzureFuncRollback),
			Description:        "", //TODO
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
			Rollback:           true,
		})
	}
	return out
}

func buildPipeline(stages []sdk.StageConfig, autoRollback bool) []sdk.PipelineStage {
	out := make([]sdk.PipelineStage, 0, len(stages)+1)
	for _, s := range stages {
		out = append(out, sdk.PipelineStage{
			Name:               s.Name,
			Index:              s.Index,
			Rollback:           false,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}

	if autoRollback {
		// we set the index of the rollback stage to the minimum index of all stages.
		minIndex := slices.MinFunc(stages, func(a, b sdk.StageConfig) int {
			return a.Index - b.Index
		}).Index

		out = append(out, sdk.PipelineStage{
			Name:               string(AzureFuncRollback),
			Index:              minIndex,
			Rollback:           true,
			Metadata:           map[string]string{},
			AvailableOperation: sdk.ManualOperationNone,
		})
	}
	return out
}
