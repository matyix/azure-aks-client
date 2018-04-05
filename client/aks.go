package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/banzaicloud/azure-aks-client/cluster"
	"github.com/banzaicloud/azure-aks-client/utils"
	banzaiTypesAzure "github.com/banzaicloud/banzai-types/components/azure"
	banzaiConstants "github.com/banzaicloud/banzai-types/constants"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

const BaseUrl = "https://management.azure.com"

type AKSClient struct {
	azureSdk *cluster.Sdk
	logger   *logrus.Logger
	clientId string
	secret   string
}

// GetAKSClient creates an *AKSClient instance with the passed credentials and default logger
func GetAKSClient(credentials *cluster.AKSCredential) (*AKSClient, error) {

	azureSdk, err := cluster.Authenticate(credentials)
	if err != nil {
		return nil, err
	}
	aksClient := &AKSClient{
		clientId: azureSdk.ServicePrincipal.ClientID,
		secret:   azureSdk.ServicePrincipal.ClientSecret,
		azureSdk: azureSdk,
		logger:   getDefaultLogger(),
	}
	if aksClient.clientId == "" {
		return nil, utils.NewErr("clientID is missing")
	}
	if aksClient.secret == "" {
		return nil, utils.NewErr("secret is missing")
	}
	return aksClient, nil
}

// With sets logger
func (a *AKSClient) With(i interface{}) {
	if a != nil {
		switch i.(type) {
		case logrus.Logger:
			logger := i.(logrus.Logger)
			a.logger = &logger
		case *logrus.Logger:
			a.logger = i.(*logrus.Logger)
		}
	}
}

// getDefaultLogger return the default logger
func getDefaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Level = logrus.InfoLevel
	logger.Formatter = new(logrus.JSONFormatter)
	return logger
}

// GetCluster gets the details of the managed cluster with a specified resource group and name.
func (a *AKSClient) GetCluster(name string, resourceGroup string) (*banzaiTypesAzure.ResponseWithValue, error) {

	a.logInfof("Start getting aks cluster: %s [%s]", name, resourceGroup)

	managedCluster, err := a.azureSdk.ManagedClusterClient.Get(context.Background(), resourceGroup, name)
	if err != nil {
		return nil, err
	}

	a.logInfof("Status code: %d", managedCluster.StatusCode)

	return &banzaiTypesAzure.ResponseWithValue{
		StatusCode: managedCluster.StatusCode,
		Value:      *convertManagedClusterToValue(&managedCluster),
	}, nil
}

// ListClusters gets a list of managed clusters in the specified subscription. The operation returns properties of each managed
// cluster.
func (a *AKSClient) ListClusters(resourceGroup string) (*banzaiTypesAzure.ListResponse, error) {
	a.logInfof("Start getting cluster list from %s resource group", resourceGroup)

	list, err := a.azureSdk.ManagedClusterClient.List(context.Background())
	if err != nil {
		return nil, err
	}

	managedClusters := list.Values()

	a.logInfo("Create response model")
	response := banzaiTypesAzure.ListResponse{StatusCode: list.Response().StatusCode, Value: banzaiTypesAzure.Values{
		Value: convertManagedClustersToValues(managedClusters),
	}}
	return &response, nil
}

// convertManagedClustersToValues returns []Value with the managed clusters properties
func convertManagedClustersToValues(managedCluster []containerservice.ManagedCluster) []banzaiTypesAzure.Value {
	var values []banzaiTypesAzure.Value
	for _, mc := range managedCluster {
		values = append(values, *convertManagedClusterToValue(&mc))
	}
	return values
}

// convertManagedClusterToValue returns Value with the ManagedCluster properties
func convertManagedClusterToValue(managedCluster *containerservice.ManagedCluster) *banzaiTypesAzure.Value {
	return &banzaiTypesAzure.Value{
		Id:         *managedCluster.ID,
		Location:   *managedCluster.Location,
		Name:       *managedCluster.Name,
		Properties: banzaiTypesAzure.Properties{},
	}
}

// CreateUpdateCluster creates or updates a managed cluster with the specified configuration for agents and Kubernetes
// version.
func (a *AKSClient) CreateUpdateCluster(request *cluster.CreateClusterRequest) (*banzaiTypesAzure.ResponseWithValue, error) {

	if request == nil {
		return nil, errors.New("Empty request")
	}

	a.logInfo("Start create/update cluster")
	a.logDebugf("CreateRequest: %v", request)
	a.logInfo("Validate cluster create/update request")

	if err := request.Validate(); err != nil {
		return nil, err
	}
	a.logInfo("Validate passed")

	managedCluster := cluster.GetManagedCluster(request, a.clientId, a.secret)
	a.logDebugf("Created managed cluster model - %#v", &managedCluster)

	a.logDebug("Send request to azure")
	result, err := a.azureSdk.ManagedClusterClient.CreateOrUpdate(context.Background(), request.ResourceGroup, request.Name, *managedCluster)
	if err != nil {
		return nil, err
	}

	resp := result.Response()

	a.logDebugf("Read response body: %v", resp.Body)
	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprint("error during cluster creation:", err)
		return nil, utils.NewErr(msg)
	}

	a.logInfof("Status code: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// something went wrong, create failed
		errResp := utils.CreateErrorFromValue(resp.StatusCode, value)
		return nil, errResp
	}

	a.logInfo("Create response model")
	v := banzaiTypesAzure.Value{}
	json.Unmarshal([]byte(value), &v)

	return &banzaiTypesAzure.ResponseWithValue{
		StatusCode: resp.StatusCode,
		Value:      v,
	}, nil
}

// DeleteCluster deletes the managed cluster with a specified resource group and name.
func (a *AKSClient) DeleteCluster(name string, resourceGroup string) error {
	a.logInfof("Start deleting cluster %s in %s resource group", name, resourceGroup)
	a.logDebug("Send request to azure")
	response, err := a.azureSdk.ManagedClusterClient.Delete(context.Background(), resourceGroup, name)
	if err != nil {
		return err
	}

	a.logInfof("Status code: %d", response.Response().StatusCode)

	return nil
}

// PollingCluster polls until the cluster ready or an error occurs
func (a *AKSClient) PollingCluster(name string, resourceGroup string) (*banzaiTypesAzure.ResponseWithValue, error) {
	const stageSuccess = "Succeeded"
	const stageFailed = "Failed"
	const waitInSeconds = 10

	a.logInfof("Start polling cluster: %s [%s]", name, resourceGroup)

	a.logDebug("Start loop")

	result := banzaiTypesAzure.ResponseWithValue{}
	for isReady := false; !isReady; {

		a.logDebug("Send request to azure")
		managedCluster, err := a.azureSdk.ManagedClusterClient.Get(context.Background(), resourceGroup, name)
		if err != nil {
			return nil, err
		}

		statusCode := managedCluster.StatusCode
		a.logInfof("Cluster polling status code: %d", statusCode)

		convertManagedClusterToValue(&managedCluster)

		switch statusCode {
		case http.StatusOK:
			response := convertManagedClusterToValue(&managedCluster)

			stage := *managedCluster.ProvisioningState
			a.logInfof("Cluster stage is %s", stage)

			switch stage {
			case stageSuccess:
				isReady = true
				result.Update(http.StatusCreated, *response)
			case stageFailed:
				return nil, banzaiConstants.ErrorAzureCLusterStageFailed
			default:
				a.logInfo("Waiting for cluster ready...")
				time.Sleep(waitInSeconds * time.Second)
			}

		default:
			return nil, errors.New("status code is not OK")
		}
	}

	return &result, nil
}

// GetClusterConfig gets the given cluster kubeconfig
func (a *AKSClient) GetClusterConfig(name, resourceGroup, roleName string) (*banzaiTypesAzure.Config, error) {

	a.logInfof("Start getting %s cluster's config in %s, role name: %s", name, resourceGroup, roleName)

	a.logDebug("Send request to azure")
	profile, err := a.azureSdk.ManagedClusterClient.GetAccessProfiles(context.Background(), resourceGroup, name, roleName)
	if err != nil {
		return nil, err
	}

	a.logInfof("Status code: %d", profile.StatusCode)
	a.logInfo("Create response model")
	return &banzaiTypesAzure.Config{
		Location: *profile.Location,
		Name:     *profile.Name,
		Properties: struct {
			KubeConfig string `json:"kubeConfig"`
		}{
			KubeConfig: string(*profile.KubeConfig),
		},
	}, nil
}

// GetVmSizes lists all available virtual machine sizes for a subscription in a location.
func (a *AKSClient) GetVmSizes(location string) ([]string, error) {

	a.logInfo("Start listing vm sizes")
	resp, err := a.azureSdk.VMSizeClient.List(context.Background(), location)
	if err != nil {
		return nil, err
	}

	var sizes []string
	for _, vm := range *resp.Value {
		sizes = append(sizes, *vm.Name)
	}
	return sizes, nil
}

// GetLocations returns all the locations that are available for resource providers
func (a *AKSClient) GetLocations() ([]string, error) {

	a.logInfo("Start listing locations")
	resp, err := a.azureSdk.SubscriptionsClient.ListLocations(context.Background(), a.azureSdk.ServicePrincipal.SubscriptionID)
	if err != nil {
		return nil, err
	}

	var locations []string
	for _, loc := range *resp.Value {
		locations = append(locations, *loc.Name)
	}

	return locations, nil
}

// GetKubernetesVersions returns a list of supported kubernetes version in the specified subscription
func (a *AKSClient) GetKubernetesVersions(location string) ([]string, error) {

	a.logInfo("Start listing Kubernetes versions")
	resp, err := a.azureSdk.ContainerServicesClient.ListOrchestrators(context.Background(), location, string(compute.Kubernetes))
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, v := range *resp.OrchestratorVersionProfileProperties.Orchestrators {
		versions = append(versions, *v.OrchestratorVersion)
	}

	return versions, nil
}

func (a *AKSClient) logDebug(args ...interface{}) {
	if a.logger != nil {
		a.logger.Debug(args...)
	}
}
func (a *AKSClient) logInfo(args ...interface{}) {
	if a.logger != nil {
		a.logger.Info(args...)
	}
}
func (a *AKSClient) logWarn(args ...interface{}) {
	if a.logger != nil {
		a.logger.Warn(args...)
	}
}
func (a *AKSClient) logError(args ...interface{}) {
	if a.logger != nil {
		a.logger.Error(args...)
	}
}

func (a *AKSClient) logFatal(args ...interface{}) {
	if a.logger != nil {
		a.logger.Fatal(args...)
	}
}

func (a *AKSClient) logPanic(args ...interface{}) {
	if a.logger != nil {
		a.logger.Panic(args...)
	}
}

func (a *AKSClient) logDebugf(format string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Debugf(format, args...)
	}
}

func (a *AKSClient) logInfof(format string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Infof(format, args...)
	}
}

func (a *AKSClient) logWarnf(format string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Warnf(format, args...)
	}
}

func (a *AKSClient) logErrorf(format string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Errorf(format, args...)
	}
}

func (a *AKSClient) logFatalf(format string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Fatalf(format, args...)
	}
}

func (a *AKSClient) logPanicf(format string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Panicf(format, args...)
	}
}
