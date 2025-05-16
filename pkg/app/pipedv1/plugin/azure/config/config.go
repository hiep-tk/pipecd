package config

type AzureDeployTargetConfig struct {
	SubscriptionID string `json:"subscriptionID"`
	//shared key file, env file, env variable?
}

type AzureApplicationSpec struct {
	Kind             string         `json:"kind"`
	FunctionManifest *FunctionsSpec `json:"functionManifest"`
}

const (
	FunctionKind = "function"
)
