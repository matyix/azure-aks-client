package client

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-06-01/subscriptions"
	"github.com/banzaicloud/azure-aks-client/cluster"
	"github.com/banzaicloud/azure-aks-client/utils"
	"github.com/sirupsen/logrus"
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

func (a *AKSClient) List() (containerservice.ManagedClusterListResultPage, error) {
	return a.azureSdk.ManagedClusterClient.List(context.Background())
}

func (a *AKSClient) CreateOrUpdate(request *cluster.CreateClusterRequest, managedCluster *containerservice.ManagedCluster) (containerservice.ManagedClustersCreateOrUpdateFuture, error) {
	return a.azureSdk.ManagedClusterClient.CreateOrUpdate(context.Background(), request.ResourceGroup, request.Name, *managedCluster)
}

func (a *AKSClient) Delete(resourceGroup, name string) (containerservice.ManagedClustersDeleteFuture, error) {
	return a.azureSdk.ManagedClusterClient.Delete(context.Background(), resourceGroup, name)
}

func (a *AKSClient) Get(resourceGroup, name string) (containerservice.ManagedCluster, error) {
	return a.azureSdk.ManagedClusterClient.Get(context.Background(), resourceGroup, name)
}

func (a *AKSClient) GetAccessProfiles(resourceGroup, name, roleName string) (containerservice.ManagedClusterAccessProfile, error) {
	return a.azureSdk.ManagedClusterClient.GetAccessProfiles(context.Background(), resourceGroup, name, roleName)
}

func (a *AKSClient) ListVmSizes(location string) (result compute.VirtualMachineSizeListResult, err error) {
	return a.azureSdk.VMSizeClient.List(context.Background(), location)
}

func (a *AKSClient) ListLocations() (subscriptions.LocationListResult, error) {
	return a.azureSdk.SubscriptionsClient.ListLocations(context.Background(), a.azureSdk.ServicePrincipal.SubscriptionID)
}

func (a *AKSClient) ListVersions(location, resourceType string) (result containerservice.OrchestratorVersionProfileListResult, err error) {
	return a.azureSdk.ContainerServicesClient.ListOrchestrators(context.Background(), location, resourceType)
}

func (a *AKSClient) GetClientId() string {
	return a.clientId
}

func (a *AKSClient) GetClientSecret() string {
	return a.secret
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
