package initapi

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/banzaicloud/azure-aks-client/cluster"
	banzaiTypes "github.com/banzaicloud/banzai-types/components"
	"github.com/banzaicloud/azure-aks-client/utils"
)

var sdk cluster.Sdk

const AzureClientId = "AZURE_CLIENT_ID"
const AzureClientSecret = "AZURE_CLIENT_SECRET"
const AzureSubscriptionId = "AZURE_SUBSCRIPTION_ID"
const AzureTenantId = "AZURE_TENANT_ID"

func Authenticate() (*cluster.Sdk, *banzaiTypes.BanzaiResponse) {
	clientId := os.Getenv(AzureClientId)
	clientSecret := os.Getenv(AzureClientSecret)
	subscriptionId := os.Getenv(AzureSubscriptionId)
	tenantId := os.Getenv(AzureTenantId)

	// ---- [Check Environmental variables] ---- //
	if len(clientId) == 0 {
		return nil, utils.CreateEnvErrorResponse(AzureClientId)
	}

	if len(clientSecret) == 0 {
		return nil, utils.CreateEnvErrorResponse(AzureClientSecret)
	}

	if len(subscriptionId) == 0 {
		return nil, utils.CreateEnvErrorResponse(AzureSubscriptionId)
	}

	if len(tenantId) == 0 {
		return nil, utils.CreateEnvErrorResponse(AzureTenantId)
	}

	sdk = cluster.Sdk{
		ServicePrincipal: &cluster.ServicePrincipal{
			ClientID:       clientId,
			ClientSecret:   clientSecret,
			SubscriptionID: subscriptionId,
			TenantId:       tenantId,
			HashMap: map[string]string{
				AzureClientId:       clientId,
				AzureClientSecret:   clientSecret,
				AzureSubscriptionId: subscriptionId,
				AzureTenantId:       tenantId,
			},
		},
	}

	authenticatedToken, err := helpers.NewServicePrincipalTokenFromCredentials(sdk.ServicePrincipal.HashMap, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, utils.CreateAuthErrorResponse(err)
	}

	sdk.ServicePrincipal.AuthenticatedToken = authenticatedToken

	resourceGroup := resources.NewGroupsClient(sdk.ServicePrincipal.SubscriptionID)
	resourceGroup.Authorizer = autorest.NewBearerAuthorizer(sdk.ServicePrincipal.AuthenticatedToken)
	sdk.ResourceGroup = &resourceGroup

	return &sdk, nil
}

func GetSdk() *cluster.Sdk {
	return &sdk
}
