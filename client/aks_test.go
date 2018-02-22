package client_test

import (
	"testing"
	"github.com/banzaicloud/azure-aks-client/client"
	"github.com/banzaicloud/azure-aks-client/cluster"
)

const name = "test_cluster"
const resourceGroup = "test_rg1"
const roleName = "clusterUser"

func TestCreateCluster(t *testing.T) {

	t.Log(" --- [ Testing creation ] ---")

	cl, err := client.GetAKSClient(nil)
	if err != nil {
		t.Errorf("Error during get aks client %v", err)
		t.FailNow() // todo probald ki
	}
	cl.With(nil) // logging out

	c := cluster.CreateClusterRequest{
		Name:              name,
		Location:          "eastus",
		VMSize:            "Standard_D2_v2",
		ResourceGroup:     resourceGroup,
		AgentCount:        1,
		AgentName:         "agentpool1",
		KubernetesVersion: "1.7.7",
	}

	if resp, err := cl.CreateUpdateCluster(c); err != nil {
		t.Errorf("Error is NOT <nil>: %s.", err)
	} else if resp.Value.Name != name {
		t.Errorf("Expected cluster name is %v but got %v.", name, resp.Value.Name)
	}

}

func TestPollingCluster(t *testing.T) {

	t.Log(" --- [ Testing polling ] ---")
	cl, err := client.GetAKSClient(nil)
	if err != nil {
		t.Errorf("Error during get aks client %v", err)
		t.FailNow() // todo probald ki
	}
	cl.With(nil) // logging out

	if _, err := cl.PollingCluster(name, resourceGroup); err != nil {
		t.Errorf("Error is NOT <nil>: %s. Polling failed.", err)
	}

}

func TestListCluster(t *testing.T) {

	t.Log(" --- [ Testing listing ] ---")
	cl, err := client.GetAKSClient(nil)
	if err != nil {
		t.Errorf("Error during get aks client %v", err)
		t.FailNow() // todo probald ki
	}
	cl.With(nil) // logging out

	if resp, err := cl.ListClusters(resourceGroup); err != nil {
		t.Errorf("Error is NOT <nil>: %s. Listing failed.", err)
	} else {

		isContains := false
		for i := 0; i < len(resp.Value.Value); i++ {
			v := resp.Value.Value[i]
			if v.Name == name {
				isContains = true
				break
			}
		}

		if !isContains {
			t.Errorf("The list not contains %v in %v", name, resourceGroup)
		}
	}

}

func TestGetK8sConfig(t *testing.T) {

	t.Log(" --- [ Testing delete ] ---")

	cl, err := client.GetAKSClient(nil)
	if err != nil {
		t.Errorf("Error during get aks client %v", err)
		t.FailNow() // todo probald ki
	}
	cl.With(nil) // logging out

	if _, err := cl.GetClusterConfig(name, resourceGroup, roleName); err != nil {
		t.Errorf("Get cluster config failed: %v", err)
		t.FailNow()
	}

}

func TestDeleteCluster(t *testing.T) {

	t.Log(" --- [ Testing delete ] ---")
	cl, err := client.GetAKSClient(nil)
	if err != nil {
		t.Errorf("Error during get aks client %v", err)
		t.FailNow() // todo probald ki
	}
	cl.With(nil) // logging out

	if err := cl.DeleteCluster(name, resourceGroup); err != nil {
		t.Errorf("Delete failed: %s.", err.Error())
	}

}
