package initapi

import (
	"github.com/banzaicloud/azure-aks-client/cluster"
)

func Init() (*cluster.Sdk, *InitErrorResponse) {
	clusterSdk, err := Authenticate()
	return clusterSdk, err
}
