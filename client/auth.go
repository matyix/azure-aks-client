package client

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/banzaicloud/azure-aks-client/utils"
	log "github.com/sirupsen/logrus"
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

var sdk Sdk

func init() {
	// Log as JSON
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func Authenticate() *resources.GroupsClient {
	clientId := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	tenantId := os.Getenv("AZURE_TENANT_ID")

	sdk = Sdk{
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

	if err := utils.CheckEnvVar(&sdk.ServicePrincipal.HashMap); err != nil {
		log.WithFields(log.Fields{
			"Environment check error": err,
		}).Error("Environment variables missing")
		return nil
	}

	authenticatedToken, err := helpers.NewServicePrincipalTokenFromCredentials(sdk.ServicePrincipal.HashMap, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.WithFields(log.Fields{
			"Authentication error": err,
		}).Error("Failed to authenticate with Azure")
		return nil
	}

	sdk.ServicePrincipal.AuthenticatedToken = authenticatedToken

	resourceGroup := resources.NewGroupsClient(sdk.ServicePrincipal.SubscriptionID)
	resourceGroup.Authorizer = autorest.NewBearerAuthorizer(sdk.ServicePrincipal.AuthenticatedToken)
	sdk.ResourceGroup = &resourceGroup

	return sdk.ResourceGroup
}

func GetSdk() *Sdk {
	return &sdk
}
