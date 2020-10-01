package serviceendpoint

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceServiceEndpointSonarQube schema and implementation for SonarQube service endpoint resource
func ResourceServiceEndpointSonarQube() *schema.Resource {
	r := genBaseServiceEndpointResource(flattenServiceEndpointSonarQube, expandServiceEndpointSonarQube)

	r.Schema["url"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotWhiteSpace,
		Description:  "Url for the SonarQube Server",
	}

	r.Schema["token"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotWhiteSpace,
		Description:  "Authentication Token generated through SonarQube (go to My Account > Security > Generate Tokens)",
	}

	r.Schema["all_pipelines"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Grant access permission to all pipelines",
	}

	return r
}

// Convert internal Terraform data structure to an AzDO data structure
func expandServiceEndpointSonarQube(d *schema.ResourceData) (*serviceendpoint.ServiceEndpoint, *string, error) {

	serviceEndpoint, projectID := doBaseExpansion(d)
	serviceEndpoint.Authorization = &serviceendpoint.EndpointAuthorization{
		Scheme: converter.String("UsernamePassword"),
		Parameters: &map[string]string{
			"username": d.Get("token").(string),
		},
	}
	serviceEndpoint.IsShared = converter.Bool(d.Get("all_pipelines").(bool))
	serviceEndpoint.Type = converter.String("sonarqube")
	serviceEndpoint.Url = converter.String(d.Get("url").(string))
	return serviceEndpoint, projectID, nil
}

// Convert AzDO data structure to internal Terraform data structure
func flattenServiceEndpointSonarQube(d *schema.ResourceData, serviceEndpoint *serviceendpoint.ServiceEndpoint, projectID *string) {

	doBaseFlattening(d, serviceEndpoint, projectID)
	d.Set("url", *serviceEndpoint.Url)
	d.Set("all_pipelines", *serviceEndpoint.IsShared)
	d.Set("token", (*serviceEndpoint.Authorization.Parameters)["username"])
}

// Example Terraform:
//
// resource "azuredevops_serviceendpoint_sonarqube" "example" {
// 		project_id = azuredevops_project.project.id
// 		url   = "https://sonarqube.instance.com"
// 		token = "xxx"
// 		service_endpoint_name = "example-sonarqube-1"
// 		description           = "example-1"
//   }

// CLI: https://docs.microsoft.com/en-us/cli/azure/ext/azure-devops/devops/service-endpoint?view=azure-cli-latest#ext_azure_devops_az_devops_service_endpoint_create
// How to obtain Config: https://docs.microsoft.com/en-us/azure/devops/cli/service-endpoint?view=azure-devops#create-service-endpoint-using-configuration-file
//
// Example API Request (via Chrome developer tools)
// {
// 	"administratorsGroup": null,
// 	"authorization": {
// 		"parameters": {
// 			"username": "redactedredactedredactedredactedredactedredacted"
// 		},
// 		"scheme": "UsernamePassword"
// 	},
// 	"createdBy": null,
// 	"data": {},
// 	"description": "My Big Fat Description",
// 	"groupScopeId": null,
// 	"isShared": false,
// 	"name": "My Big Fat Service Connection Name",
// 	"operationStatus": null,
// 	"owner": "library",
// 	"readersGroup": null,
// 	"serviceEndpointProjectReferences": [
// 	{
// 		"description": "My Big Fat Description",
// 		"name": "My Big Fat Service Connection Name",
// 		"projectReference": {
// 		"id": "deadbeef-cafe-cafe-deadbeef-deadbeef",
// 		"name": "My Big Fat Project"
// 		}
// 	}
// 	],
// 	"type": "sonarqube",
// 	"url": "http://example.com"
// }
