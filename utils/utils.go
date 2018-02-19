package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	banzaiConstants "github.com/banzaicloud/banzai-types/constants"
	banzaiUtils "github.com/banzaicloud/banzai-types/utils"
	banzaiTypes "github.com/banzaicloud/banzai-types/components"
)

const (
	credentialsPath = "/.azure/credentials.json"
)

// ToJSON returns the passed item as a pretty-printed JSON string. If any JSON error occurs,
// it returns the empty string.
func ToJSON(v interface{}) (string, error) {
	j, err := json.MarshalIndent(v, "", "  ")
	return string(j), err
}

// NewServicePrincipalTokenFromCredentials creates a new ServicePrincipalToken using values of the
// passed credentials map.
func NewServicePrincipalTokenFromCredentials(c map[string]string, scope string) (*adal.ServicePrincipalToken, error) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, c["AZURE_TENANT_ID"])
	if err != nil {
		panic(err)
	}
	return adal.NewServicePrincipalToken(*oauthConfig, c["AZURE_CLIENT_ID"], c["AZURE_CLIENT_SECRET"], scope)
}

func ensureValueStrings(mapOfInterface map[string]interface{}) map[string]string {
	mapOfStrings := make(map[string]string)
	for key, value := range mapOfInterface {
		mapOfStrings[key] = ensureValueString(value)
	}
	return mapOfStrings
}

func ensureValueString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func ReadPubRSA(filename string) string {
	b, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/" + filename)
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

func CheckEnvVar(envVars *map[string]string) error {
	var missingVars []string
	for varName, value := range *envVars {
		if value == "" {
			missingVars = append(missingVars, varName)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("Missing environment variables %v", missingVars)
	}
	return nil
}

func S(input string) *string {
	s := input
	return &s
}

type AzureServerError struct {
	Message string `json:"message"`
}

func CreateErrorFromValue(statusCode int, v []byte) AzureServerError {
	if statusCode == banzaiConstants.BadRequest {
		ase := AzureServerError{}
		json.Unmarshal([]byte(v), &ase)
		if len(ase.Message) != 0 {
			return ase
		}
	}

	type TempError struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	tempError := TempError{}
	json.Unmarshal([]byte(v), &tempError)
	return AzureServerError{Message: tempError.Error.Message}
}

func CreateEnvErrorResponse(env string) *banzaiTypes.BanzaiResponse {
	message := "Environmental variable is empty: " + env
	banzaiUtils.LogError(banzaiConstants.TagInit, "environmental_error")
	return &banzaiTypes.BanzaiResponse{StatusCode: banzaiConstants.InternalErrorCode, Message: message}
}

func CreateAuthErrorResponse(err error) *banzaiTypes.BanzaiResponse {
	errMsg := "Failed to authenticate with Azure"
	banzaiUtils.LogError(banzaiConstants.TagAuth, "Authentication error:", err)
	return &banzaiTypes.BanzaiResponse{StatusCode: banzaiConstants.InternalErrorCode, Message: errMsg}
}
