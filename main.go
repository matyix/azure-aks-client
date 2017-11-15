package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

type Sdk struct {
	ServicePrincipal *ServicePrincipal
	ResourceGroup    *resources.GroupsClient
}

type ServicePrincipal struct {
	ClientID           string
	ClientSecret       string
	SubscriptionID     string
	TenantId           string
	HashMap            map[string]string
	AuthenticatedToken *adal.ServicePrincipalToken
}

func main() {
	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		fmt.Errorf("Empty $AZURE_CLIENT_ID")
		return
	}
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if clientSecret == "" {
		fmt.Errorf("Empty $AZURE_CLIENT_SECRET")
		return
	}
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		fmt.Errorf("Empty $AZURE_SUBSCRIPTION_ID")
		return
	}
	tenantID := os.Getenv("AZURE_TENANT_ID")
	if tenantID == "" {

		return
	}

	sdk := &Sdk{
		ServicePrincipal: &ServicePrincipal{
			ClientID:       clientID,
			ClientSecret:   clientSecret,
			SubscriptionID: subscriptionID,
			TenantId:       tenantID,
			HashMap: map[string]string{
				"AZURE_CLIENT_ID":       clientID,
				"AZURE_CLIENT_SECRET":   clientSecret,
				"AZURE_SUBSCRIPTION_ID": subscriptionID,
				"AZURE_TENANT_ID":       tenantID,
			},
		},
	}

	authenticatedToken, err := helpers.NewServicePrincipalTokenFromCredentials(sdk.ServicePrincipal.HashMap, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		fmt.Errorf("Token %#v", authenticatedToken)
		return
	}

	fmt.Printf("Token %#v", authenticatedToken.Token)
	sdk.ServicePrincipal.AuthenticatedToken = authenticatedToken

	resourceGroup := resources.NewGroupsClient(sdk.ServicePrincipal.SubscriptionID)
	resourceGroup.Authorizer = autorest.NewBearerAuthorizer(sdk.ServicePrincipal.AuthenticatedToken)
	sdk.ResourceGroup = &resourceGroup

	fmt.Printf("Resource Group: %v#", sdk.ResourceGroup)

	location := "eastus"
	sdk.ResourceGroup.CreateOrUpdate("myRG1", resources.Group{Location: &location})

	return
}
