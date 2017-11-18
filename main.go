package main

import (
	azure "github.com/matyix/azure-aks-client/client"
)

func main() {

	cluster := azure.ClusterDetails{
		Name:          "AK47-reloaded",
		ResourceGroup: "rg1",
		Location:      "eastus",
		VMSize:        "Standard_D2_v2",
		DNSPrefix:     "gun",
		AdminUsername: "",
		PubKeyName:    "id_rsa.pub",
	}

	//azure.ListClusters(azure.Authenticate())
	azure.CreateCluster(azure.Authenticate(), cluster)
	//azure.DeleteCluster(azure.Authenticate(), cluster)
}
