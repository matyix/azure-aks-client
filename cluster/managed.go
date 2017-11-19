package cluster

import (
	"github.com/Azure/azure-sdk-for-go/arm/containerservice"
	utils "github.com/banzaicloud/azure-aks-client/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

}

type CreateRequest struct {
	Location   string            `json:"location"`
	Properties ClusterProperties `json:"properties"`
}

type ClusterProperties struct {
	DNSPrefix               string                                   `json:"dnsPrefix"`
	AgentPoolProfiles       []AgentPoolProfiles                      `json:"agentPoolProfiles"`
	KubernetesVersion       string                                   `json:"kubernetesVersion"`
	LinuxProfile            containerservice.LinuxProfile            `json:"linuxProfile"`
	ServicePrincipalProfile containerservice.ServicePrincipalProfile `json:"servicePrincipalProfile"`
}

type AgentPoolProfiles struct {
	Count  int    `json:"count"`
	Name   string `json:"name"`
	VMSize string `json:"vmSize"`
}

type ManagedCluster struct {
	ClusterDetails    ClusterDetails
	ClusterProperties ClusterProperties
}

type ClusterDetails struct {
	Name          string
	ResourceGroup string
	Location      string
	VMSize        string
	DNSPrefix     string
	AdminUsername string
	PubKeyName    string
}

func GetManagedCluster() *ManagedCluster {
	return &ManagedCluster{
		ClusterDetails: ClusterDetails{
			Name:          "AK47-reloaded",
			ResourceGroup: "rg1",
			Location:      "eastus",
			VMSize:        "Standard_D2_v2",
			DNSPrefix:     "gun",
			AdminUsername: "",
			PubKeyName:    "id_rsa.pub",
		},
		ClusterProperties: ClusterProperties{
			DNSPrefix: clusterDetails.DNSPrefix,
			AgentPoolProfiles: []AgentPoolProfiles{
				{
					Count:  1,
					Name:   "agentpool1",
					VMSize: clusterDetails.VMSize,
				},
			},
			KubernetesVersion: "1.7.7",
			ServicePrincipalProfile: containerservice.ServicePrincipalProfile{
				//ClientID: &clientId,
				//Secret:   &clientSecret,
				ClientID: &"cid",
				Secret:   &"secret",
			},
			LinuxProfile: containerservice.LinuxProfile{
				//AdminUsername: S(clusterDetails.AdminUsername),
				AdminUsername: S("ubuntu"),
				SSH: &containerservice.SSHConfiguration{
					PublicKeys: &[]containerservice.SSHPublicKey{
						{
							//KeyData: utils.S(utils.ReadPubRSA(clusterDetails.PubKeyName)),
							KeyData: utils.S(utils.ReadPubRSA("id_rsa")),
						},
					},
				},
			},
		},
	}
}
