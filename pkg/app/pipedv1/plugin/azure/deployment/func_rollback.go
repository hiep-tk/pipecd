package deployment

import (
	"context"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/azure/config"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
)

func (p *Plugin) executeAzureFuncRollbackStage(ctx context.Context, dts []*sdk.DeployTarget[config.AzureDeployTargetConfig], input *sdk.ExecuteStageInput[config.AzureApplicationSpec]) sdk.StageStatus {
	return sdk.StageStatusSuccess
}
