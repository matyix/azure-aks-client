package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	//Create resource group
	/**
	location := "eastus"
	sdk.ResourceGroup.CreateOrUpdate("myRG1", resources.Group{Location: &location})
	**/

	//Autorest test
	/**
		GET /subscriptions/SUBSCRIPTION_ID/resourcegroups?api-version=2015-01-01 HTTP/1.1
		Host: management.azure.com
		Authorization: Bearer YOUR_ACCESS_TOKEN
		Content-Type: application/json
	**/

	p := map[string]interface{}{"subscription-id": subscriptionID}
	q := map[string]interface{}{"api-version": "2015-01-01"}

	req, _ := autorest.Prepare(&http.Request{},
		resourceGroup.WithAuthorization(),
		autorest.AsGet(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourcegroups", p),
		autorest.WithQueryParameters(q))

	resp, err := autorest.SendWithSender(resourceGroup.Client, req)
	if err != nil {
		fmt.Errorf("SendWithSender %#v", err)
		return
	}

	value := struct {
		ResourceGroups []struct {
			Name string `json:"name"`
		} `json:"value"`
	}{}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&value)
	if err != nil {
		fmt.Errorf("Decode %#v", err)
		return
	}

	var groupNames = make([]string, len(value.ResourceGroups))
	fmt.Printf("Groups : %#v", groupNames)
	for i, name := range value.ResourceGroups {
		groupNames[i] = name.Name
	}

	fmt.Println("Groups:", strings.Join(groupNames, ", "))

	return

}
