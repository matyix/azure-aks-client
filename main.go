package main

import (
	azure "github.com/banzaicloud/azure-aks-client/client"
	"os"
)

func main() {

	azure.ListClusters(azure.Authenticate(), os.Getenv("AZURE_SUBSCRIPTION_ID"))

}
