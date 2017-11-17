package client

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/arm/containerservice"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/matyix/azure-aks-client/utils"
	"io/ioutil"
	"net/http"
	"os"
)

type AgentPoolProfiles struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	VMSize string `json:"vmSize"`
}

type ClusterProperties struct {
	DNSPrefix               string                                   `json:"dnsPrefix"`
	AgentPoolProfiles       []AgentPoolProfiles                      `json:"agentPoolProfiles"`
	KubernetesVersion       string                                   `json:"kubernetesVersion"`
	LinuxProfile            containerservice.LinuxProfile            `json:"linuxProfile"`
	ServicePrincipalProfile containerservice.ServicePrincipalProfile `json:"servicePrincipalProfile"`
}

type CreateRequestBody struct {
	Location   string            `json:"location"`
	Properties ClusterProperties `json:"properties"`
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

	fmt.Printf("Resp status : %#v\n", resp.StatusCode)
	value, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("ListResponse: %#v\n", string(value))

	respListInGR := ListInRG{}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&respListInGR)
	if err != nil {
		fmt.Errorf("Decode %#v", err)
		return
	}

	fmt.Printf("List in RG : %#v", &respListInGR)

	fmt.Printf("Resp status : %#v\n", resp.StatusCode)

}

/*
CreateCluster creates a managed AKS on Azure
PUT https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31
*/

func S(input string) *string {
	s := input
	return &s
}
func CreateCluster(groupClient *resources.GroupsClient, subscriptionId, name string) {

	pathParam := map[string]interface{}{"subscription-id": subscriptionId, "resourceName": name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}
	clientId := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	createRequest := CreateRequestBody{
		Location: "eastus",
		Properties: ClusterProperties{
			DNSPrefix: "dnsprefix1",
			AgentPoolProfiles: []AgentPoolProfiles{
				{
					Count:  1,
					Name:   "agentpool1",
					VMSize: "Standard_D2_v2",
				},
			},
			KubernetesVersion: "1.7.7",
			ServicePrincipalProfile: containerservice.ServicePrincipalProfile{
				ClientID: &clientId,
				Secret:   &clientSecret,
			},
			LinuxProfile: containerservice.LinuxProfile{
				AdminUsername: S("admin"),
				SSH: &containerservice.SSHConfiguration{
					PublicKeys: &[]containerservice.SSHPublicKey{
						{
							KeyData: S(utils.ReadPubRSA("id_rsa.pub")),
						},
					},
				},
			},
		},
	}
	//if clusterProperties != nil {
	//	createRequest.properties = clusterProperties
	//}

	req, _ := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsPut(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/rg1/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
		autorest.WithJSON(createRequest),
		autorest.AsContentType("application/json"),
	)

	val, err := json.Marshal(createRequest)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("JSONbody: %s", val)

	resp, err := autorest.SendWithSender(groupClient.Client, req)

	defer resp.Body.Close()
	//dec := json.NewDecoder(resp.Body)
	//err = dec.Decode(&ListInRG{})
	if err != nil {
		fmt.Errorf("Decode %#v", err)
		return
	}

	fmt.Printf("Resp status : %#v", resp.StatusCode)
	value, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%#v", string(value))

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
