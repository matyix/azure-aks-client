package main

import (
	client "github.com/banzaicloud/azure-aks-client/client"
	cluster "github.com/banzaicloud/azure-aks-client/cluster"
)

func main() {

	cluster := cluster.GetManagedCluster()

	//azure.ListClusters(azure.Authenticate())
	client.CreateCluster(client.Authenticate(), cluster)
	//azure.DeleteCluster(azure.Authenticate(), cluster)
}
