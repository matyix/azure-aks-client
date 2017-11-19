package main

import (
	client "github.com/banzaicloud/azure-aks-client/client"
	cluster "github.com/banzaicloud/azure-aks-client/cluster"
)

func main() {

	var sdk cluster.Sdk
	sdk = *client.Authenticate()

	clientId := sdk.ServicePrincipal.ClientID
	secret := sdk.ServicePrincipal.ClientSecret

	cluster := cluster.GetTestManagedCluster(clientId, secret)

	//azure.ListClusters(azure.Authenticate())
	client.CreateCluster(&sdk, *cluster, "lofasz", "rg1")
	//azure.DeleteCluster(azure.Authenticate(), cluster)
}
