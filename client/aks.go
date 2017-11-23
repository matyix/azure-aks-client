package client

import (
	"encoding/json"
	"github.com/Azure/go-autorest/autorest"
	"github.com/banzaicloud/azure-aks-client/cluster"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func init() {
	// Log as JSON
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

const internalErrorCode = 500

/*
ListClusters is listing AKS clusters in the specified subscription and resource group
GET https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters?
	api-version=2017-08-31
*/
func ListClusters(sdk *cluster.Sdk, resourceGroup string, initError *InitErrorResponse) string {

	if sdk == nil {
		return initError.toString()
	}

	pathParam := map[string]interface{}{
		"subscription-id": sdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   resourceGroup}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *sdk.ResourceGroup

	req, err := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsGet(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters", pathParam),
		autorest.WithQueryParameters(queryParam))

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during listing clusters in ", resourceGroup, " resource group")
		return createErrorResponse()
	}

	log.Info("Start cluster listing in ", resourceGroup, " resource group")

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during listing clusters in ", resourceGroup, " resource group")
		return createErrorResponse()
	}

	log.Info("Cluster list response status code: ", resp.StatusCode)

	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during listing clusters in ", resourceGroup, " resource group")
		return createErrorResponse()
	}

	azureListResponse := AzureListResponse{}
	json.Unmarshal([]byte(value), &azureListResponse)
	log.Info("List cluster result ", azureListResponse.toString())

	response := ListResponse{StatusCode: resp.StatusCode, Value: azureListResponse}
	return response.toString()
}

/*
CreateCluster creates a managed AKS on Azure
PUT https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31sdk *cluster.Sdk
*/
func CreateCluster(sdk *cluster.Sdk, managedCluster cluster.ManagedCluster, name string, resourceGroup string, initError *InitErrorResponse) string {

	if sdk == nil {
		return initError.toString()
	}

	pathParam := map[string]interface{}{
		"subscription-id": sdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   resourceGroup,
		"resourceName":    name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *sdk.ResourceGroup

	req, _ := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsPut(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
		autorest.WithJSON(managedCluster),
		autorest.AsContentType("application/json"),
	)

	_, err := json.Marshal(managedCluster)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during JSON marshal")
		return createErrorResponse()
	}

	log.Info("Cluster creation start with name ", name, " in ", resourceGroup, " resource group")

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during cluster creation")
		return createErrorResponse()
	}

	defer resp.Body.Close()
	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during cluster creation")
		return createErrorResponse()
	}

	log.Info("Cluster create response code: ", resp.StatusCode)

	response := Value{}
	json.Unmarshal([]byte(value), &response)
	log.Info("Cluster creation with name ", name, " in ", resourceGroup, " resource group has started")

	result := Response{StatusCode: resp.StatusCode, Value: response}
	return result.toString()
}

/*
DeleteCluster deletes a managed AKS on Azure
DELETE https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31
*/
func DeleteCluster(sdk *cluster.Sdk, name string, resourceGroup string, initError *InitErrorResponse) string {

	if sdk == nil {
		return initError.toString()
	}

	pathParam := map[string]interface{}{
		"subscription-id": sdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   resourceGroup,
		"resourceName":    name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *sdk.ResourceGroup

	req, err := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsDelete(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
	)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during delete cluster")
		return createErrorResponse()
	}

	log.Info("Cluster delete start with name ", name, " in ", resourceGroup, " resource group")

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during delete cluster")
		return createErrorResponse()
	}

	log.Info("Cluster delete status code: ", resp.StatusCode)

	defer resp.Body.Close()
	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during delete cluster")
		return createErrorResponse()
	}

	log.Info("Delete cluster response ", string(value))

	result := ResponseWithCode{StatusCode: resp.StatusCode}
	return result.toString()
}

/*
PollingCluster polling AKS on Azure
GET https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31
 */
func PollingCluster(sdk *cluster.Sdk, name string, resourceGroup string, initError *InitErrorResponse) string {

	if sdk == nil {
		return initError.toString()
	}

	const OK = 200
	const stageSuccess = "Succeeded"
	const stageFailed = "Failed"
	const waitInSeconds = 10

	pathParam := map[string]interface{}{
		"subscription-id": sdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   resourceGroup,
		"resourceName":    name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *sdk.ResourceGroup

	req, err := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsGet(),
		autorest.WithBaseURL("https://management.azure.com"),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
	)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during cluster polling")
		return createErrorResponse()
	}

	log.Info("Cluster polling start with name ", name, " in ", resourceGroup, " resource group")

	result := Response{}
	for isReady := false; !isReady; {

		resp, err := autorest.SendWithSender(groupClient.Client, req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("error during cluster polling")
			return createErrorResponse()
		}

		statusCode := resp.StatusCode
		log.Info("Cluster polling status code: ", statusCode)

		switch statusCode {
		case OK:
			value, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.WithFields(log.Fields{"error": err}).Error("error during cluster polling")
				return createErrorResponse()
			}

			response := Value{}
			json.Unmarshal([]byte(value), &response)

			stage := response.Properties.ProvisioningState
			log.Info("Cluster stage is ", stage)

			switch stage {
			case stageSuccess:
				isReady = true
				result.update(statusCode, response)
			case stageFailed:
				return createErrorResponse()
			default:
				log.Info("Waiting...")
				time.Sleep(waitInSeconds * time.Second)
			}

		default:
			return createErrorResponseWithCode(statusCode)
		}
	}

	return result.toString()
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

type AzureListResponse struct {
	Value []Value `json:"value"`
}

type Value struct {
	Id         string     `json:"id"`
	Location   string     `json:"location"`
	Name       string     `json:"name"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	ProvisioningState string    `json:"provisioningState"`
	AgentPoolProfiles []Profile `json:"agentPoolProfiles"`
}

type Profile struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Response struct {
	StatusCode int   `json:"status_code"`
	Value      Value `json:"message"`
}

type ListResponse struct {
	StatusCode int               `json:"status_code"`
	Value      AzureListResponse `json:"message"`
}

type ResponseWithCode struct {
	StatusCode int    `json:"status_code"`
}

func (r AzureListResponse) toString() string {
	jsonResponse, _ := json.Marshal(r)
	return string(jsonResponse)
}

func (v Value) toString() string {
	jsonResponse, _ := json.Marshal(v)
	return string(jsonResponse)
}

func (r Response) toString() string {
	jsonResponse, _ := json.Marshal(r)
	return string(jsonResponse)
}

func (r *Response) update(code int, Value Value) {
	r.Value = Value
	r.StatusCode = code
}

func createErrorResponse() string {
	return createErrorResponseWithCode(internalErrorCode)
}

func createErrorResponseWithCode(code int) string {
	errorResponse := ResponseWithCode{StatusCode: code}
	return errorResponse.toString()
}

func (r ListResponse) toString() string {
	jsonResponse, _ := json.Marshal(r)
	return string(jsonResponse)
}

func (r ResponseWithCode) toString() string {
	jsonResponse, _ := json.Marshal(r)
	return string(jsonResponse)
}
