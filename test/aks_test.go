package test

import (
	"testing"
	"github.com/banzaicloud/azure-aks-client/client"
	"github.com/banzaicloud/azure-aks-client/cluster"
	"fmt"
	"encoding/json"
)

const name = "test_cluster"
const resourceGroup = "rg1"

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
	}

	result := client.Response{}
	value := client.CreateUpdateCluster(c)
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
	value := client.PollingCluster(name, resourceGroup)
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
	value := client.ListClusters(resourceGroup)
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
	value := client.DeleteCluster(name, resourceGroup)
	json.Unmarshal([]byte(value), &result)

	if result.StatusCode != 202 {
		t.Errorf("Expected response status code is 202 but got %v", result.StatusCode)
	}

}
