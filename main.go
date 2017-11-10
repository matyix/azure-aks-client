package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/go-resty/resty"
	"github.com/matyix/artisan-aks-client/utils"
)

func main() {
	resourceGroup := "BanzaiCloud"
	resourceName := "TestClster"

	c := map[string]string{
		"AZURE_CLIENT_ID":       os.Getenv("AZURE_CLIENT_ID"),
		"AZURE_CLIENT_SECRET":   os.Getenv("AZURE_CLIENT_SECRET"),
		"AZURE_SUBSCRIPTION_ID": os.Getenv("AZURE_SUBSCRIPTION_ID"),
		"AZURE_TENANT_ID":       os.Getenv("AZURE_TENANT_ID")}
	if err := checkEnvVar(&c); err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	spt, err := utils.NewServicePrincipalTokenFromCredentials(c, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	//GET
	//List clusters - GET https://management.azure.com/subscriptions/
	// {subscriptionId}/resourceGroups/
	// {resourceGroupName}/providers/Microsoft.ContainerService/managedClusters?
	// api-version=2017-08-31

	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"api-version": "2017-08-31",
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken(spt.Token.AccessToken).
		Get("https://management.azure.com/subscriptions/" +
			c["AZURE_SUBSCRIPTION_ID"] +
			"/resourceGroups/rg1/providers/Microsoft.ContainerService/managedClusters")

	printOutput(resp, err)

	type ManagedClusterProperties struct {
		accessProfiles    string
		fqdn              string
		kubernetesVersion string
		provisioningState string
	}

	type Error struct {
		/* variables */
	}
	//PUT
	//Create cluster - PUT https://management.azure.com/subscriptions/
	// {subscriptionId}/resourceGroups/
	// {resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	// api-version=2017-08-31
	resp, err = resty.R().
		SetQueryParams(map[string]string{
			"api-version": "2017-08-31",
		}).
		SetBody(ManagedClusterProperties{
			fqdn:              "banzai.com",
			kubernetesVersion: "1.8.3",
			accessProfiles:    "access",
			provisioningState: "state"}).
		SetAuthToken(spt.Token.AccessToken).
		SetError(&Error{}). // or SetError(Error{}).
		Put("https://management.azure.com/subscriptions/" + c["AZURE_SUBSCRIPTION_ID"] +
			"/resourceGroups/" + resourceGroup +
			"/providers/Microsoft.ContainerService/managedClusters/" + resourceName)

	printOutput(resp, err)

}

func printOutput(resp *resty.Response, err error) {
	fmt.Println(resp, err)
}

func checkEnvVar(envVars *map[string]string) error {
	var missingVars []string
	for varName, value := range *envVars {
		if value == "" {
			missingVars = append(missingVars, varName)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("Missing environment variables %v", missingVars)
	}
	return nil
}
