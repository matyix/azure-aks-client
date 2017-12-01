package cluster

import (
	"github.com/banzaicloud/azure-aks-client/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log as JSON
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

}

type ManagedCluster struct {
	Location   string     `json:"location"`
	Properties Properties `json:"properties"`
}

func GetManagedCluster(request CreateClusterRequest) *ManagedCluster {
	return &ManagedCluster{
		Location: request.Location,
		Properties: Properties{
			DNSPrefix: "dnsprefix",
			AgentPoolProfiles: []AgentPoolProfiles{
				{
					Count:  request.AgentCount,
					Name:   request.AgentName,
					VMSize: request.VMSize,
				},
			},
			KubernetesVersion: "1.7.7",
			ServicePrincipalProfile: ServicePrincipalProfile{
				ClientID: utils.S(request.ClientId),
				Secret:   utils.S(request.Secret),
			},
			LinuxProfile: LinuxProfile{
				AdminUsername: "erospista",
				SSH: SSH{
					PublicKeys: &[]SSHPublicKey{
						{
							KeyData: utils.S(utils.ReadPubRSA("id_rsa.pub")),
						},
					},
				},
			},
		},
	}
}

type CreateClusterRequest struct {
	Name              string
	Location          string
	VMSize            string
	ResourceGroup     string
	AgentCount        int
	AgentName         string
	KubernetesVersion string
	ClientId          string
	Secret            string
}
