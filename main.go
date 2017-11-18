package main

import (
	azure "github.com/matyix/azure-aks-client/client"
)

func main() {

	cluster := azure.ClusterDetails{
		Name:          "AK47-reloaded",
		Location:      "eastus",
		VMSize:        "Standard_D2_v2",
		DNSPrefix:     "gun",
		AdminUsername: "faszacsavo1",
		PubKeyName:    "id_rsa.pub",
	}

	//azure.ListClusters(azure.Authenticate(), os.Getenv("AZURE_SUBSCRIPTION_ID"))
	azure.CreateCluster(azure.Authenticate(), cluster)
	//azure.DeleteCluster(azure.Authenticate(), os.Getenv("AZURE_SUBSCRIPTION_ID"), "AK47-reloaded")
}
