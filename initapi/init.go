package initapi

import (
	"github.com/banzaicloud/azure-aks-client/client"
	"github.com/banzaicloud/azure-aks-client/cluster"
)

func Init() (*cluster.Sdk, *client.InitErrorResponse) {
	clusterSdk, err := client.Authenticate()
	return clusterSdk, err
}
