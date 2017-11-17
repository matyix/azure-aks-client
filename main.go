package main

import (
	azure "github.com/matyix/azure-aks-client/client"
	"os"
)

func main() {

	//azure.ListClusters(azure.Authenticate(), os.Getenv("AZURE_SUBSCRIPTION_ID"))
	azure.CreateCluster(azure.Authenticate(), os.Getenv("AZURE_SUBSCRIPTION_ID"), "AK-47-S")

}
