package main

import (
	"fmt"
	"github.com/banzaicloud/azure-aks-client/client"
	"github.com/banzaicloud/azure-aks-client/cluster"
	"github.com/banzaicloud/azure-aks-client/initapi"
)

var sdk *cluster.Sdk
var initError *client.InitErrorResponse

func init() {
	sdk, initError = initapi.Init()
}

func main() {

	clientId := ""
	secret := ""
	if sdk != nil {
		clientId = sdk.ServicePrincipal.ClientID
		secret = sdk.ServicePrincipal.ClientSecret
	}

	cluster := cluster.GetTestManagedCluster(clientId, secret)
	fmt.Printf("Cluster :#%v \n", cluster)

	result := client.ListClusters(sdk, "rg1", initError)
	// result := client.CreateCluster(sdk, *cluster, "lofasz", "rg1", initError)
	// result := client.DeleteCluster(sdk, "lofasz", "rg1", initError)
	// result := client.PollingCluster(sdk,"lofasz","rg1", initError)
	fmt.Println("--------", result, "--------")
}
