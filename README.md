### AKS client until Swagger/SDK is released

This is a temporary solution until the following PR' are not merged upstream.

**For API Specification**: Go client is missing AKS until this [PR](https://github.com/Azure/azure-rest-api-specs/pull/1956) is merged in the API specification.
Related [PR](https://github.com/Azure/azure-rest-api-specs/pull/1912) superseeded by the previos one.


**For API Go client**: The API itself is lacking the AKS feature until this [issue](https://github.com/Azure/azure-sdk-for-go/issues/847) is fixed.

This is a library to create Microsoft Managed Kubernetes clusters (**AKS**) on Azure cloud.

#### Prerequisities 

You will need the following ENV variables exported: `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`, `AZURE_TENANT_ID`, `AZURE_SUBSCRIPTION_ID`

You can get this information from the portal, but the easiest and fastes way is to use the Azure CLI tool.

Install the tool and login using the following commands.

```bash
$ curl -L https://aka.ms/InstallAzureCli | bash
$ exec -l $SHELL
$ az login
```

Create a `Service Principal` for the Azure Active Directory using the following command.

```bash
$ az ad sp create-for-rbac

```

You should get something like: 

``` 
{

  "appId": "1234567-1234-1234-1234-1234567890ab",
  "displayName": "azure-cli-2017-08-18-19-25-59",
  "name": "http://azure-cli-2017-08-18-19-25-59",
  "password": "1234567-1234-1234-be18-1234567890ab",
  "tenant": "1234567-1234-1234-be18-1234567890ab"
}
```

Translate the output from the previous command to newly exported environmental variables.

Service Principal Variable Name | Environmental variable
--- | ---
appId | AZURE_CLIENT_ID
password | AZURE_CLIENT_SECRET
tenant | AZURE_TENANT_ID

Run the following command to get you Azure subscription ID.

```bash
$ az account show --query id
"1234567-1234-1234-1234567890ab"
```

Finally export that value as an environmental variable as well.

Command| Environmental variable
--- | ---
az account show --query id | AZURE_SUBSCRIPTION_ID

**At this point you should have the following 4 environmental variables set!**

```bash
export AZURE_CLIENT_ID = "1234567-1234-1234-1234567890ab"
export AZURE_CLIENT_SECRET = "1234567-1234-1234-1234567890ab"
export AZURE_TENANT_ID = "1234567-1234-1234-1234567890ab"
export AZURE_SUBSCRIPTION_ID = "1234567-1234-1234-1234567890ab"
```