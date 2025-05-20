package deployment

import (
	"context"
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/azure/config"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/azure/provider"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"

	"strings"
)

func (p *Plugin) executeAzureFuncSyncStage(ctx context.Context, dts []*sdk.DeployTarget[config.AzureDeployTargetConfig], input *sdk.ExecuteStageInput[config.AzureApplicationSpec]) sdk.StageStatus {
	lp := input.Client.LogPersister()
	lp.Info("Start azure func sync")
	if len(dts) != 1 {
		lp.Errorf("Currently support only one deployment target, instead got %d", len(dts))
		return sdk.StageStatusFailure
	}
	sdkClient, err := provider.NewAzureClient(ctx, dts[0].Config, lp)
	if err != nil {
		lp.Errorf("Failed to create Azure Functions sdkClient: %v", err)
		return sdk.StageStatusFailure
	}
	sdkClient.SetResourceTags(map[string]*string{
		provider.LabelManagedBy:   to.Ptr(provider.ManagedByPiped),
		provider.LabelPiped:       to.Ptr(input.Request.Deployment.PipedID),
		provider.LabelCommitHash:  to.Ptr(input.Request.TargetDeploymentSource.CommitHash),
		provider.LabelApplication: to.Ptr(input.Request.Deployment.ApplicationID),
	})
	appCfg, err := input.Request.TargetDeploymentSource.AppConfig()
	if err != nil {
		lp.Errorf("Failed to get AppConfig: %v", err)
		return sdk.StageStatusFailure
	}
	manifest := appCfg.Spec.FunctionManifest
	if manifest == nil {
		lp.Errorf("AppConfig.Spec.FunctionManifest is nil")
		return sdk.StageStatusFailure
	}
	stageConfigDump := input.Request.StageConfig
	var slotNames []string
	if len(stageConfigDump) > 0 {
		var stageConfig config.AzureFunctionSyncStageConfig
		if err = json.Unmarshal(stageConfigDump, &stageConfig); err != nil {
			lp.Errorf("cannot read sync stage config: %v", err)
		}
		slotNames = append(slotNames, stageConfig.SlotName)
	}
	if manifest.ArmTemplate != nil {
		lp.Infof("Start using arm template to deploy %s: template %s, parameter %s", manifest.ArmTemplate.DeploymentName, manifest.ArmTemplate.DeploymentTemplateFile, manifest.ArmTemplate.DeploymentParameterFile)
		err = sdkClient.DeployARMTemplate(ctx, manifest.ResourceGroupName, input.Request.TargetDeploymentSource.ApplicationDirectory, *manifest.ArmTemplate)
		if err != nil {
			lp.Errorf("Failed to deploy ARM template: %v", err)
			return sdk.StageStatusFailure
		}
	}
	needCreated, err := sdkClient.FunctionValidate(ctx, manifest.ResourceGroupName, manifest.FunctionName, slotNames)
	if err != nil {
		lp.Errorf("Failed to validate manifest: %v", err)
		return sdk.StageStatusFailure
	}
	if needCreated {
		lp.Errorf("Cannot find resource even after deploy template")
		return sdk.StageStatusFailure
	}
	var slotName string
	if len(slotNames) > 0 {
		slotName = slotNames[0]
	}
	current, err := sdkClient.FunctionGetX(ctx, manifest.ResourceGroupName, manifest.FunctionName, slotName)
	if err != nil {
		lp.Errorf("Failed to get Function %s: %v", manifest.FunctionName, err)
		return sdk.StageStatusFailure
	}

	if strings.Contains(*current.Kind, "linux") && *current.Properties.SKU == "Dynamic" { //Consumption Linux plan special treatment
		if err = sdkClient.FunctionRunFromPackageDeploymentX(ctx, manifest.ResourceGroupName, manifest.FunctionName, slotName, manifest.PackageUri); err != nil {
			lp.Errorf("Failed to deploy %s with 'WEBSITE_RUN_FROM_PACKAGE': %v", manifest.FunctionName, err)
			return sdk.StageStatusFailure
		}
	} else {
		if err = sdkClient.FunctionKuduDeploymentX(ctx, manifest.ResourceGroupName, manifest.FunctionName, slotName, manifest.PackageUri); err != nil {
			lp.Errorf("Failed to deploy %s with KuduDeployment: %v", manifest.FunctionName, err)
			return sdk.StageStatusFailure
		}
	}
	return sdk.StageStatusSuccess
}
