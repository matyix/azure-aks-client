package client

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

func init() {
	// Log as JSON
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

/*
ListClusters is listing AKS clusters in the specified subscription and resource group
GET https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters?
	api-version=2017-08-31
*/
func ListClusters(groupClient *resources.GroupsClient, managedCluster ManagedCluster) {

	pathParam := map[string]interface{}{
		"subscription-id": sdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   managedCluster.ClusterDetails.ResourceGroup}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	req, _ := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsGet(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters", pathParam),
		autorest.WithQueryParameters(queryParam))

	resp, err := autorest.SendWithSender(groupClient.Client, req)

	log.Info("REST call response status %#v", resp.StatusCode)
	value, err := ioutil.ReadAll(resp.Body)
	log.Info("Cluster list response", string(value))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error during cluster list call ")
		return
	}

	respListInGR := ListInRG{}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&respListInGR)
	/*
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error during cluster list decode ")
			return
		}
	*/

	log.Info("List cluster call response status", resp.StatusCode)
	log.Info("Cluster list in the resource group", &respListInGR)

}

/*
CreateCluster creates a managed AKS on Azure
PUT https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31
*/
func CreateCluster(groupClient *resources.GroupsClient, managedCluster ManagedCluster) {

	pathParam := map[string]interface{}{
		"subscription-id": sdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   managedCluster.ClusterDetails.ResourceGroup,
		"resourceName":    managedCluster.ClusterDetails.Name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	createRequest := CreateRequest{
		Location:   managedCluster.ClusterDetails.Location,
		Properties: managedCluster.ClusterProperties,
	}
	//if clusterProperties != nil {
	//	createRequest.properties = clusterProperties
	//}

	req, _ := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsPut(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
		autorest.WithJSON(createRequest),
		autorest.AsContentType("application/json"),
	)

	//val, err := json.Marshal(createRequest)
	_, err := json.Marshal(createRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error during JSON marshal ")
		return
	}
	//log.Info("JSON body ", val)

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error during cluster create call ")
		return
	}

	defer resp.Body.Close()
	value, err := ioutil.ReadAll(resp.Body)
	/*
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error during cluster create decode ")
			return
		}
	*/
	log.Info("Cluster create call response status", resp.StatusCode)
	log.Info("Cluster create response", string(value))

}

/*
DeleteCluster deletes a managed AKS on Azure
DELETE https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31
*/
func DeleteCluster(groupClient *resources.GroupsClient, managedCluster ManagedCluster) {

	pathParam := map[string]interface{}{
		"subscription-id": sdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   clusterDetails.ResourceGroup,
		"resourceName":    clusterDetails.Name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	req, _ := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsDelete(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
	)

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error during cluster delete call ")
		return
	}

	log.Info("Delete cluster call response status", resp.StatusCode)
	value, err := ioutil.ReadAll(resp.Body)
	log.Info("delete cluster response", string(value))

	respListInGR := ListInRG{}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&respListInGR)
	/*
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error during cluster delete decode ")
			return
		}
	*/
	log.Info("Delete cluster call response status", resp.StatusCode)

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
