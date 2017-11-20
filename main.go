package main

import (
	"fmt"
	client "github.com/banzaicloud/azure-aks-client/client"
	cluster "github.com/banzaicloud/azure-aks-client/cluster"
	"github.com/banzaicloud/azure-aks-client/initapi"
)

var sdk *cluster.Sdk

func init() {
	sdk = initapi.Init()
}

func main() {

	clientId := sdk.ServicePrincipal.ClientID
	secret := sdk.ServicePrincipal.ClientSecret

	cluster := cluster.GetTestManagedCluster(clientId, secret)
	fmt.Printf("Cluster :#%v ", cluster)

	//client.ListClusters(&sdk, "rg1")
	client.CreateCluster(sdk, *cluster, "lofasz", "rg1")
	//client.DeleteCluster(sdk, "lofasz", "rg1")
}
