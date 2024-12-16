package e2etests

import (
	"fmt"
	"os"
	"testing"

	"github.com/kong/candidate-take-home-exercise-sdet/internal/config"
	"github.com/kong/candidate-take-home-exercise-sdet/internal/server"
	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
	services "github.com/kong/candidate-take-home-exercise-sdet/test/services"
	"github.com/stretchr/testify/assert"
)

var (
	auth          *services.TokensService
	client        *framework.HttpClient
	configuration config.Config
)

// TestMain sets up and tears down global resources before and after the tests.
func TestMain(m *testing.M) {
	// Setup before tests
	fmt.Println("Setting up tests...")
	baseUrl := "https://dummyjson.com"
	client = framework.NewHttpClient(baseUrl)
	auth = services.NewTokensService(client, "http://localhost:18080")
	configuration = framework.GetConfiguration()

	// Run the tests
	code := m.Run() // This runs all the tests in the package.

	// Teardown after tests
	fmt.Println("Tearing down tests...")
	// Clean up resources if needed (e.g., closing connections, deleting temp files, etc.)
	// You can also do service cleanup here if needed.

	// Exit with the code returned by m.Run()
	// This allows Go to handle the exit code of the tests.
	os.Exit(code)
}

func TestAuthService_CreateToken_InvalidUsername(t *testing.T) {

	payload := CreateCustomCredentialsReqBody(configuration.Username+"INVALID", configuration.Password)
	resp, err := auth.CreateToken(payload)
	assert.Equal(t, 401, resp.StatusCode)
	assert.Nil(t, err.Error)
	errorBody, _ := framework.ParseResponseBody[models.ErrorResponse](resp.Body)
	assert.Equal(t, "Invalid username or password", errorBody.Error)
}

func TestAuthService_CreateToken_InvalidPassword(t *testing.T) {

	payload := CreateCustomCredentialsReqBody(configuration.Username, configuration.Password+"INVALID")
	resp, err := auth.CreateToken(payload)
	assert.Equal(t, 401, resp.StatusCode)
	assert.Nil(t, err.Error)
	errorBody, _ := framework.ParseResponseBody[models.ErrorResponse](resp.Body)
	assert.Contains(t, errorBody.Error, "password is not equal to ")
}

func TestAuthService_CreateToken_CheckTokenValidity(t *testing.T) {

	payload := CreateCredentialsReqBody()
	resp, err := auth.CreateToken(payload)
	assert.Nil(t, err.Error)
	token, _ := framework.ParseResponseBody[server.TokenResponse](resp.Body)
	assert.True(t, len(token.Token) > 5)
	tokenisValid, _ := framework.TokenHasUsernameClaim(token.Token, configuration.Username)
	assert.True(t, tokenisValid)

}

func CreateCredentialsReqBody() server.Credentials {
	creds := server.Credentials{Username: configuration.Username, Password: configuration.Password}
	return creds
}

func CreateCustomCredentialsReqBody(username string, password string) server.Credentials {
	creds := server.Credentials{Username: username, Password: password}
	return creds
}
