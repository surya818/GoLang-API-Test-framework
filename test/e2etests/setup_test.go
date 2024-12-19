package e2etests

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/kong/candidate-take-home-exercise-sdet/internal/config"
	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
	service "github.com/kong/candidate-take-home-exercise-sdet/test/services"
)

var (
	Client            *framework.HttpClient
	AuthorizationApi  *service.TokensService
	Configuration     config.Config
	ServiceApi        *service.ServiceApi
	ServiceVersionApi *service.ServiceVersionApi
	token             string
)

// Initialize sets up common test configurations.
func Initialize() {
	fmt.Println("Setting up tests...")

	// Common setup
	baseUrl := "http://localhost:18080"
	Client = framework.NewHttpClient(baseUrl)
	AuthorizationApi = service.NewTokensService(Client, baseUrl)
	Configuration = framework.GetConfiguration()
	token = GetToken()
	ServiceApi = service.NewServiceApi(Client, baseUrl, token)
	ServiceVersionApi = service.NewServiceVersionApi(Client, baseUrl, token)
	err := framework.InitLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
	}
}

// Teardown cleans up resources after tests.
func Teardown() {
	framework.Logger.Info("Tests finished") // Add any teardown logic here if needed.
}

// RunWithSetup initializes resources and runs tests for the e2etests package.
func RunWithSetup(m *testing.M) {
	Initialize()
	code := m.Run() // Run tests
	Teardown()
	os.Exit(code) // Exit with test code
}

func TestMain(m *testing.M) {
	RunWithSetup(m) // Call the shared setup and teardown logic
}

func CreateService(payload models.Service) (http.Response, framework.ApiError) {
	service_resp, service_err := ServiceApi.CreateService(payload)
	return service_resp, service_err
}

func extractServiceResponse(service_resp http.Response) models.ServiceResponse {
	resp_object, _ := framework.ParseResponseBody[models.ServiceResponse](service_resp.Body)
	return resp_object

}

func extractErrorResponse(service_resp http.Response) models.ErrorResponse {
	resp_object, _ := framework.ParseResponseBody[models.ErrorResponse](service_resp.Body)
	return resp_object

}

func extractServiceVersionResponse(service_version_resp http.Response) models.ServiceVersionResponse {
	resp_object, _ := framework.ParseResponseBody[models.ServiceVersionResponse](service_version_resp.Body)
	return resp_object

}

func extractListServicesResponse(service_resp http.Response) models.ListServices {
	resp_object, _ := framework.ParseResponseBody[models.ListServices](service_resp.Body)
	return resp_object

}

func extractListServiceVersionsResponse(service_resp http.Response) models.ListServiceVersions {
	resp_object, _ := framework.ParseResponseBody[models.ListServiceVersions](service_resp.Body)
	return resp_object

}

func listServicesAndExtractTheList() models.ListServices {
	listServices, _ := ServiceApi.ListServices()
	services := extractListServicesResponse(listServices)
	return services
}

func listServiceVersionsAndExtractTheList(serviceId string) models.ListServiceVersions {
	listServiceVersions, _ := ServiceVersionApi.ListServiceVersions(serviceId)
	serviceVersionsList := extractListServiceVersionsResponse(listServiceVersions)
	return serviceVersionsList
}

func serviceWithIDExists(services models.ListServices, serviceId string) (bool, models.Service) {
	var found bool
	var targetService models.Service
	for _, service := range services.Items {
		if service.ID == serviceId {
			found = true
			targetService = service

			break
		}
	}
	return found, targetService
}

func serviceVersionWithIDExists(serviceVersions models.ListServiceVersions, versionId string) (bool, models.ServiceVersion) {
	var found bool
	var targetService models.ServiceVersion
	for _, serviceVersion := range serviceVersions.Items {
		if serviceVersion.ID == versionId {
			found = true
			targetService = serviceVersion

			break
		}
	}
	return found, targetService
}

func GetToken() string {

	authToken := os.Getenv("AUTH_TOKEN")

	// Check if it is empty
	if authToken == "" {
		authToken = AuthorizationApi.FetchToken(Configuration.Username, Configuration.Password)
		_ = os.Setenv("AUTH_TOKEN", authToken) // Optionally set it in the environment
	} else {
		framework.Logger.Info("Auth Token extracted locally")

	}

	return authToken
}

func CreateService_Success() models.ServiceResponse {
	serviceName := framework.GetRandomName("service")
	payload := framework.CreateServicePayload(serviceName, serviceName, "test service")
	service_resp, service_err := CreateService(payload)
	fmt.Errorf("Error in creating Service: %v", service_err)
	if service_resp.StatusCode != 201 {
		fmt.Errorf("Error in creating Service: Status code is %v", service_resp.StatusCode)
		return models.ServiceResponse{}
	}
	service_object := extractServiceResponse(service_resp)
	return service_object
}

func CreateServiceVersion_Success() models.ServiceVersionResponse {
	service_object := CreateService_Success()
	serviceId := service_object.Item.ID
	payload := framework.CreateServiceVersionPayload(serviceId, "", "")
	service_version_resp, service_err := ServiceVersionApi.CreateServiceVersion(serviceId, payload)
	fmt.Errorf("Error in creating Service: %v", service_err)
	if service_version_resp.StatusCode != 201 {
		fmt.Errorf("Error in creating Service: Status code is %v", service_version_resp.StatusCode)
		return models.ServiceVersionResponse{}
	}
	service_version_object := extractServiceVersionResponse(service_version_resp)
	return service_version_object
}
