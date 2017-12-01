package test

import (
	"testing"
	"github.com/banzaicloud/azure-aks-client/client"
	"github.com/banzaicloud/azure-aks-client/initapi"
	"github.com/banzaicloud/azure-aks-client/cluster"
	"fmt"
	"encoding/json"
)

var sdk *cluster.Sdk
var initError *client.InitErrorResponse
var clientId string
var secret string

const name = "test_lofasz"
const resourceGroup = "rg1"

func init() {
	sdk, initError = initapi.Init()
	if sdk != nil {
		clientId = sdk.ServicePrincipal.ClientID
		secret = sdk.ServicePrincipal.ClientSecret
	}
}

func TestCreateCluster(t *testing.T) {

	fmt.Println(" --- [ Testing creation ] ---")

	c := cluster.CreateClusterRequest{
		Name:              name,
		Location:          "eastus",
		VMSize:            "Standard_D2_v2",
		ResourceGroup:     resourceGroup,
		AgentCount:        1,
		AgentName:         "agentpool1",
		KubernetesVersion: "1.7.7",
		ClientId:          clientId,
		Secret:            secret,
	}

	result := client.Response{}
	value := client.CreateCluster(sdk, c, initError)
	json.Unmarshal([]byte (value), &result)

	if result.StatusCode != 200 && result.StatusCode != 201 {
		t.Errorf("Expected response status code is 201 or 200 but got %v", result.StatusCode)
	}

	if result.Value.Name != name {
		t.Errorf("Expected cluster name is %v but got %v", name, result.Value.Name)
	}

}

func TestPollingCluster(t *testing.T) {

	fmt.Println(" --- [ Testing polling ] ---")

	result := client.Response{}
	value := client.PollingCluster(sdk, name, resourceGroup, initError)
	json.Unmarshal([]byte (value), &result)

	if result.StatusCode != 200 {
		t.Errorf("Expected response status code is 200 but got %v", result.StatusCode)
	}

	if result.Value.Name != name {
		t.Errorf("Expected name is %v but got %v", name, result.Value.Name)
	}

}

func TestListCluster(t *testing.T) {

	fmt.Println(" --- [ Testing listing ] ---")

	result := client.ListResponse{}
	value := client.ListClusters(sdk, resourceGroup, initError)
	json.Unmarshal([]byte(value), &result)

	if result.StatusCode != 200 {
		t.Errorf("Expected response status code is 200 but got %v", result.StatusCode)
	}

	isContains := false
	for i := 0; i < len(result.Value.Value); i++ {
		v := result.Value.Value[i]
		if v.Name == name {
			isContains = true
			break
		}
	}

	if !isContains {
		t.Errorf("The list not contains %v in %v", name, resourceGroup)
	}
}

func TestDeleteCluster(t *testing.T) {

	fmt.Println(" --- [ Testing delete ] ---")

	result := client.Response{}
	value := client.DeleteCluster(sdk, name, resourceGroup, initError)
	json.Unmarshal([]byte(value), &result)

	if result.StatusCode != 202 {
		t.Errorf("Expected response status code is 202 but got %v", result.StatusCode)
	}

}
