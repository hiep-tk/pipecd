package config

type AzureDeployTemplate struct {
	DeploymentName          string `json:"deploymentName"`
	DeploymentTemplateFile  string `json:"deploymentTemplateFile"`
	DeploymentParameterFile string `json:"deploymentParameterFile"`
}
type FunctionsSpec struct {
	FunctionName      string               `json:"functionName"`
	ResourceGroupName string               `json:"resourceGroupName"`
	ArmTemplate       *AzureDeployTemplate `json:"armTemplate"`
	BicepTemplate     *AzureDeployTemplate `json:"bicepTemplate"`
	PackageUri        string               `json:"packageUri"`
}

type AzureFunctionSyncStageConfig struct {
	SlotName string `json:"slot"`
}
type AzureFunctionSwapStageConfig struct {
	SlotName1 string `json:"slot1"`
	SlotName2 string `json:"slot2"`
}
