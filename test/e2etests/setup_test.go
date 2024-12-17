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
	"github.com/stretchr/testify/assert"
)

var (
	Client           *framework.HttpClient
	AuthorizationApi *service.TokensService
	Configuration    config.Config
	ServiceApi       *service.ServiceApi
)

// Initialize sets up common test configurations.
func Initialize() {
	fmt.Println("Setting up tests...")

	// Common setup
	baseUrl := "http://localhost:18080"
	Client = framework.NewHttpClient(baseUrl)
	AuthorizationApi = service.NewTokensService(Client, baseUrl)
	Configuration = framework.GetConfiguration()
	ServiceApi = service.NewServiceApi(Client, baseUrl, "")
	ServiceApi.AuthToken = GetToken()
}

// Teardown cleans up resources after tests.
func Teardown() {
	fmt.Println("Tearing down tests...")
	// Add any teardown logic here if needed.
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

func CreateServiceAndExtractResponse(payload models.Service, t *testing.T) models.ServiceResponse {
	service_resp, service_err := CreateService(payload)
	assert.Equal(t, 201, service_resp.StatusCode)
	assert.Nil(t, service_err.Error)
	return extractServiceResponse(service_resp)
}

func CreateService(payload models.Service) (http.Response, framework.ApiError) {
	service_resp, service_err := ServiceApi.CreateService(payload)
	return service_resp, service_err
}

func extractServiceResponse(service_resp http.Response) models.ServiceResponse {
	resp_object, _ := framework.ParseResponseBody[models.ServiceResponse](service_resp.Body)
	return resp_object

}

func extractListServicesResponse(service_resp http.Response) models.ListServices {
	resp_object, _ := framework.ParseResponseBody[models.ListServices](service_resp.Body)
	return resp_object

}

func listServicesAndExtractTheList() models.ListServices {
	listServices, _ := ServiceApi.ListServices()
	services := extractListServicesResponse(listServices)
	return services
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

func GetToken() string {

	authToken := os.Getenv("AUTH_TOKEN")
	fmt.Println("Auth Token extracted locally")

	// Check if it is empty
	if authToken == "" {
		authToken = AuthorizationApi.FetchToken(Configuration.Username, Configuration.Password)
		_ = os.Setenv("AUTH_TOKEN", authToken) // Optionally set it in the environment
	}

	return authToken
}
