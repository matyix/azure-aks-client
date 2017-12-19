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
	"github.com/banzaicloud/azure-aks-client/initapi"
	"errors"
	banzaiTypes "github.com/banzaicloud/banzai-types/components"
	banzaiTypesAzure "github.com/banzaicloud/banzai-types/components/azure"
)

func init() {
	// Log as JSON
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	azureSdk, initError = initapi.Init()
	if azureSdk != nil {
		clientId = azureSdk.ServicePrincipal.ClientID
		secret = azureSdk.ServicePrincipal.ClientSecret
	}
}

const BaseUrl = "https://management.azure.com"

var azureSdk *cluster.Sdk
var clientId string
var secret string
var initError *banzaiTypes.BanzaiResponse

/**
GetCluster gets the details of the managed cluster with a specified resource group and name.
GET https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/
	{resourceName}?api-version=2017-08-31
 */
func GetCluster(name string, resourceGroup string) (*banzaiTypesAzure.ResponseWithValue, *banzaiTypes.BanzaiResponse) {

	if azureSdk == nil {
		return nil, initError
	}

	if len(clientId) == 0 || len(secret) == 0 {
		message := "ClientId or secret is empty"
		log.WithFields(log.Fields{"error": "environmental_error"}).Error(message)
		return nil, &banzaiTypes.BanzaiResponse{StatusCode: initapi.InternalErrorCode, Message: message}
	}

	pathParam := map[string]interface{}{
		"subscription-id": azureSdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   resourceGroup,
		"resourceName":    name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *azureSdk.ResourceGroup

	req, err := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsGet(),
		autorest.WithBaseURL(BaseUrl),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam))

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during listing clusters in ", resourceGroup, " resource group")
		return nil, createErrorResponseFromError(err)
	}

	log.Info("Get cluster details with name: ", name, " in ", resourceGroup, " resource group")

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during listing clusters in ", resourceGroup, " resource group")
		return nil, createErrorResponseFromError(err)
	}

	log.Info("Get Cluster response status code: ", resp.StatusCode)

	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during get cluster in ", resourceGroup, " resource group")
		return nil, createErrorResponseFromError(err)
	}

	if resp.StatusCode != initapi.OK {
		// not ok, probably 404
		errResp := initapi.CreateErrorFromValue(resp.StatusCode, value)
		log.Info("Get cluster failed with message: ", errResp.Message)
		return nil, &banzaiTypes.BanzaiResponse{StatusCode: resp.StatusCode, Message: errResp.Message}
	} else {
		// everything is ok
		v := banzaiTypesAzure.Value{}
		json.Unmarshal([]byte(value), &v)
		response := banzaiTypesAzure.ResponseWithValue{}
		response.Update(resp.StatusCode, v)
		return &response, nil
	}

}

/*
ListClusters is listing AKS clusters in the specified subscription and resource group
GET https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters?
	api-version=2017-08-31
*/
func ListClusters(resourceGroup string) (*banzaiTypesAzure.ListResponse, *banzaiTypes.BanzaiResponse) {

	if azureSdk == nil {
		return nil, initError
	}

	if len(clientId) == 0 || len(secret) == 0 {
		message := "ClientId or secret is empty"
		log.WithFields(log.Fields{"error": "environmental_error"}).Error(message)
		return nil, &banzaiTypes.BanzaiResponse{StatusCode: initapi.InternalErrorCode, Message: message}
	}

	pathParam := map[string]interface{}{
		"subscription-id": azureSdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   resourceGroup}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *azureSdk.ResourceGroup

	req, err := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsGet(),
		autorest.WithBaseURL(BaseUrl),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters", pathParam),
		autorest.WithQueryParameters(queryParam))

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during listing clusters in ", resourceGroup, " resource group")
		return nil, createErrorResponseFromError(err)
	}

	log.Info("Start cluster listing in ", resourceGroup, " resource group")

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during listing clusters in ", resourceGroup, " resource group")
		return nil, createErrorResponseFromError(err)
	}

	log.Info("Cluster list response status code: ", resp.StatusCode)

	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during listing clusters in ", resourceGroup, " resource group")
		return nil, createErrorResponseFromError(err)
	}

	if resp.StatusCode != initapi.OK {
		// not ok, probably 404
		errResp := initapi.CreateErrorFromValue(resp.StatusCode, value)
		log.Info("Listing clusters failed with message: ", errResp.Message)
		return nil, &banzaiTypes.BanzaiResponse{StatusCode: resp.StatusCode, Message: errResp.Message}
	}

	azureListResponse := banzaiTypesAzure.Values{}
	json.Unmarshal([]byte(value), &azureListResponse)
	log.Info("List cluster result ", &azureListResponse)

	response := banzaiTypesAzure.ListResponse{StatusCode: resp.StatusCode, Value: azureListResponse}
	return &response, nil
}

/*
CreateUpdateCluster creates or updates a managed cluster
PUT https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31sdk *cluster.Sdk
*/
func CreateUpdateCluster(request cluster.CreateClusterRequest) (*banzaiTypesAzure.ResponseWithValue, *banzaiTypes.BanzaiResponse) {

	if azureSdk == nil {
		return nil, initError
	}

	if len(clientId) == 0 || len(secret) == 0 {
		message := "ClientId or secret is empty"
		log.WithFields(log.Fields{"error": "environmental_error"}).Error(message)
		return nil, &banzaiTypes.BanzaiResponse{StatusCode: initapi.InternalErrorCode, Message: message}
	}

	if isValid, errMsg := request.Validate(); !isValid {
		return nil, &banzaiTypes.BanzaiResponse{StatusCode: initapi.BadRequest, Message: errMsg}
	}

	managedCluster := cluster.GetManagedCluster(request, clientId, secret)

	pathParam := map[string]interface{}{
		"subscription-id": azureSdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   request.ResourceGroup,
		"resourceName":    request.Name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *azureSdk.ResourceGroup

	req, _ := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsPut(),
		autorest.WithBaseURL(BaseUrl),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
		autorest.WithJSON(managedCluster),
		autorest.AsContentType("application/json"),
	)

	_, err := json.Marshal(managedCluster)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during JSON marshal")
		return nil, createErrorResponseFromError(err)
	}

	log.Info("Cluster creation start with name ", request.Name, " in ", request.ResourceGroup, " resource group")

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during cluster creation")
		return nil, createErrorResponseFromError(err)
	}

	defer resp.Body.Close()
	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during cluster creation")
		return nil, createErrorResponseFromError(err)
	}

	log.Info("Cluster create response code: ", resp.StatusCode)

	if resp.StatusCode != initapi.OK && resp.StatusCode != initapi.Created {
		// something went wrong, create failed
		errResp := initapi.CreateErrorFromValue(resp.StatusCode, value)
		log.Info("Cluster creation failed with message: ", errResp.Message)
		return nil, &banzaiTypes.BanzaiResponse{StatusCode: resp.StatusCode, Message: errResp.Message}
	}

	v := banzaiTypesAzure.Value{}
	json.Unmarshal([]byte(value), &v)
	log.Info("Cluster creation with name ", request.Name, " in ", request.ResourceGroup, " resource group has started")

	result := banzaiTypesAzure.ResponseWithValue{StatusCode: resp.StatusCode, Value: v}
	return &result, nil
}

/*
DeleteCluster deletes a managed AKS on Azure
DELETE https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31
*/
func DeleteCluster(name string, resourceGroup string) (*banzaiTypes.BanzaiResponse) {

	if azureSdk == nil {
		return initError
	}

	if len(clientId) == 0 || len(secret) == 0 {
		message := "ClientId or secret is empty"
		log.WithFields(log.Fields{"error": "environmental_error"}).Error(message)
		return &banzaiTypes.BanzaiResponse{StatusCode: initapi.InternalErrorCode, Message: message}
	}

	pathParam := map[string]interface{}{
		"subscription-id": azureSdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   resourceGroup,
		"resourceName":    name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *azureSdk.ResourceGroup

	req, err := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsDelete(),
		autorest.WithBaseURL(BaseUrl),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
	)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during delete cluster")
		return createErrorResponseFromError(err)
	}

	log.Info("Cluster delete start with name ", name, " in ", resourceGroup, " resource group")

	resp, err := autorest.SendWithSender(groupClient.Client, req)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during delete cluster")
		return createErrorResponseFromError(err)
	}

	log.Info("Cluster delete status code: ", resp.StatusCode)

	defer resp.Body.Close()
	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during delete cluster")
		return createErrorResponseFromError(err)
	}

	if resp.StatusCode != initapi.OK && resp.StatusCode != initapi.NoContent && resp.StatusCode != initapi.Accepted {
		errResp := initapi.CreateErrorFromValue(resp.StatusCode, value)
		log.Info("Delete cluster failed with message: ", errResp.Message)
		return &banzaiTypes.BanzaiResponse{StatusCode: resp.StatusCode, Message: errResp.Message}
	}

	log.Info("Delete cluster response ", string(value))

	result := banzaiTypes.BanzaiResponse{StatusCode: resp.StatusCode}
	return &result
}

/*
PollingCluster polling AKS on Azure
GET https://management.azure.com/subscriptions/
	{subscriptionId}/resourceGroups/
	{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}?
	api-version=2017-08-31
 */
func PollingCluster(name string, resourceGroup string) (*banzaiTypesAzure.ResponseWithValue, *banzaiTypes.BanzaiResponse) {

	if azureSdk == nil {
		return nil, initError
	}

	if len(clientId) == 0 || len(secret) == 0 {
		message := "ClientId or secret is empty"
		log.WithFields(log.Fields{"error": "environmental_error"}).Error(message)
		return nil, &banzaiTypes.BanzaiResponse{StatusCode: initapi.InternalErrorCode, Message: message}
	}

	const OK = 200
	const stageSuccess = "Succeeded"
	const stageFailed = "Failed"
	const waitInSeconds = 10

	pathParam := map[string]interface{}{
		"subscription-id": azureSdk.ServicePrincipal.SubscriptionID,
		"resourceGroup":   resourceGroup,
		"resourceName":    name}
	queryParam := map[string]interface{}{"api-version": "2017-08-31"}

	groupClient := *azureSdk.ResourceGroup

	req, err := autorest.Prepare(&http.Request{},
		groupClient.WithAuthorization(),
		autorest.AsGet(),
		autorest.WithBaseURL(BaseUrl),
		autorest.WithPathParameters("/subscriptions/{subscription-id}/resourceGroups/{resourceGroup}/providers/Microsoft.ContainerService/managedClusters/{resourceName}", pathParam),
		autorest.WithQueryParameters(queryParam),
	)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("error during cluster polling")
		return nil, createErrorResponseFromError(err)
	}

	log.Info("Cluster polling start with name ", name, " in ", resourceGroup, " resource group")

	result := banzaiTypesAzure.ResponseWithValue{}
	for isReady := false; !isReady; {

		resp, err := autorest.SendWithSender(groupClient.Client, req)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("error during cluster polling")
			return nil, createErrorResponseFromError(err)
		}

		statusCode := resp.StatusCode
		log.Info("Cluster polling status code: ", statusCode)

		value, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("error during cluster polling")
			return nil, createErrorResponseFromError(err)
		}

		switch statusCode {
		case OK:
			response := banzaiTypesAzure.Value{}
			json.Unmarshal([]byte(value), &response)

			stage := response.Properties.ProvisioningState
			log.Info("Cluster stage is ", stage)

			switch stage {
			case stageSuccess:
				isReady = true
				result.Update(statusCode, response)
			case stageFailed:
				return nil, createErrorResponseFromError(errors.New("cluster stage is 'Failed'"))
			default:
				log.Info("Waiting...")
				time.Sleep(waitInSeconds * time.Second)
			}

		default:
			errResp := initapi.CreateErrorFromValue(resp.StatusCode, value)
			log.Info("Delete cluster failed with message: ", errResp.Message)
			return nil, &banzaiTypes.BanzaiResponse{StatusCode: resp.StatusCode, Message: errResp.Message}
		}
	}

	return &result, nil
}

func createErrorResponseFromError(err error) *banzaiTypes.BanzaiResponse {
	return &banzaiTypes.BanzaiResponse{
		StatusCode: initapi.InternalErrorCode,
		Message:    err.Error(),
	}
}
