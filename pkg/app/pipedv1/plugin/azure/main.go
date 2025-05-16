package main

import (
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/azure/deployment"
	"github.com/pipe-cd/pipecd/pkg/app/pipedv1/plugin/azure/livestate"
	"github.com/pipe-cd/pipecd/pkg/plugin/sdk"
	"log"
)

func main() {
	plugin, err := sdk.NewPlugin(
		"azure", "0.0.1",
		sdk.WithDeploymentPlugin(&deployment.Plugin{}),
		sdk.WithLivestatePlugin(&livestate.Plugin{}),
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := plugin.Run(); err != nil {
		log.Fatalln(err)
	}
}
