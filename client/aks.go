package client

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"net/http"
)

type ManagedClusterProperties struct {
	accessProfiles    string
	fqdn              string
	kubernetesVersion string
	provisioningState string
}

/*
ListClusters is listing AKS clusters in the specified subscription and resource group
GET https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters?
	api-version=2017-08-31
*/
func ListClusters(groupClient *resources.GroupsClient, subscriptionId string) {

	pathParam := map[string]interface{}{"subscription-id": subscriptionId}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	req, _ := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsGet(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/rg1/providers/Microsoft.ContainerService/managedClusters", pathParam),
		autorest.WithQueryParameters(queryParam))

	resp, err := autorest.SendWithSender(groupClient.Client, req)

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ListInRG{})
	if err != nil {
		fmt.Errorf("Decode %#v", err)
		return
	}

	fmt.Printf("List in RG : %#v", &ListInRG{})

}

/*
CreateCluster creates a managed AKS on Azure
PUT https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31
*/
func CreateCluster(groupClient *resources.GroupsClient, subscriptionId string, name, ap, fqdn, k8sV, ps string) {

	pathParam := map[string]interface{}{"subscription-id": subscriptionId, "resourceName": name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	clusterProperties := ManagedClusterProperties{
		accessProfiles:    ap,
		fqdn:              fqdn,
		kubernetesVersion: k8sV,
		provisioningState: ps,
	}

	req, _ := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsPut(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/rg1/providers/Microsoft.ContainerService/managedClusters", pathParam),
		autorest.WithQueryParameters(queryParam),
		autorest.WithJSON(clusterProperties))

	resp, err := autorest.SendWithSender(groupClient.Client, req)

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ListInRG{})
	if err != nil {
		fmt.Errorf("Decode %#v", err)
		return
	}

	fmt.Printf("List in RG : %#v", &ListInRG{})

}


type ListInRG struct {
	Value []struct {
		ID         string `json:"id"`
		Location   string `json:"location"`
		Name       string `json:"name"`
		Properties struct {
			AccessProfiles struct {
				ClusterAdmin struct {
					KubeConfig string `json:"kubeConfig"`
				} `json:"clusterAdmin"`
				ClusterUser struct {
					KubeConfig string `json:"kubeConfig"`
				} `json:"clusterUser"`
			} `json:"accessProfiles"`
			AgentPoolProfiles []struct {
				Count          int    `json:"count"`
				DNSPrefix      string `json:"dnsPrefix"`
				Fqdn           string `json:"fqdn"`
				Name           string `json:"name"`
				OsType         string `json:"osType"`
				Ports          []int  `json:"ports"`
				StorageProfile string `json:"storageProfile"`
				VMSize         string `json:"vmSize"`
			} `json:"agentPoolProfiles"`
			DNSPrefix         string `json:"dnsPrefix"`
			Fqdn              string `json:"fqdn"`
			KubernetesVersion string `json:"kubernetesVersion"`
			LinuxProfile      struct {
				AdminUsername string `json:"adminUsername"`
				SSH           struct {
					PublicKeys []struct {
						KeyData string `json:"keyData"`
					} `json:"publicKeys"`
				} `json:"ssh"`
			} `json:"linuxProfile"`
			ProvisioningState       string `json:"provisioningState"`
			ServicePrincipalProfile struct {
				ClientID          string      `json:"clientId"`
				KeyVaultSecretRef interface{} `json:"keyVaultSecretRef"`
				Secret            string      `json:"secret"`
			} `json:"servicePrincipalProfile"`
		} `json:"properties"`
		Tags string `json:"tags"`
		Type string `json:"type"`
	} `json:"value"`
}
