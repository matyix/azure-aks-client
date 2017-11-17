package client

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/matyix/azure-aks-client/utils"
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

func Authenticate() *resources.GroupsClient {
	clientId := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	tenantId := os.Getenv("AZURE_TENANT_ID")

	sdk := &Sdk{
		ServicePrincipal: &ServicePrincipal{
			ClientID:       clientId,
			ClientSecret:   clientSecret,
			SubscriptionID: subscriptionId,
			TenantId:       tenantId,
			HashMap: map[string]string{
				"AZURE_CLIENT_ID":       clientId,
				"AZURE_CLIENT_SECRET":   clientSecret,
				"AZURE_SUBSCRIPTION_ID": subscriptionId,
				"AZURE_TENANT_ID":       tenantId,
			},
		},
	}

	if err := utils.CheckEnvVar((&sdk.ServicePrincipal.HashMap)); err != nil {
		fmt.Errorf("Error: %v", err)
		return nil
	}

	authenticatedToken, err := helpers.NewServicePrincipalTokenFromCredentials(sdk.ServicePrincipal.HashMap, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		fmt.Errorf("Token %#v", authenticatedToken)
		return nil
	}

	fmt.Printf("Token %#v", authenticatedToken.Token)

	sdk.ServicePrincipal.AuthenticatedToken = authenticatedToken

	resourceGroup := resources.NewGroupsClient(sdk.ServicePrincipal.SubscriptionID)
	resourceGroup.Authorizer = autorest.NewBearerAuthorizer(sdk.ServicePrincipal.AuthenticatedToken)
	sdk.ResourceGroup = &resourceGroup

	return sdk.ResourceGroup

}
