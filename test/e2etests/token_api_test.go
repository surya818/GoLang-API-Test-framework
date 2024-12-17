package e2etests

import (
	"testing"

	"github.com/kong/candidate-take-home-exercise-sdet/internal/server"
	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_CreateToken_InvalidUsername(t *testing.T) {

	payload := framework.CreateCustomCredentialsReqBody(Configuration.Username+"INVALID", Configuration.Password)
	resp, err := AuthorizationApi.CreateToken(payload)
	assert.Equal(t, 401, resp.StatusCode)
	assert.Nil(t, err.Error)
	errorBody, _ := framework.ParseResponseBody[models.ErrorResponse](resp.Body)
	assert.Equal(t, "Invalid username or password", errorBody.Error)
}

func TestAuthService_CreateToken_InvalidPassword(t *testing.T) {

	payload := framework.CreateCustomCredentialsReqBody(Configuration.Username, Configuration.Password+"INVALID")
	resp, err := AuthorizationApi.CreateToken(payload)
	assert.Equal(t, 401, resp.StatusCode)
	assert.Nil(t, err.Error)
	errorBody, _ := framework.ParseResponseBody[models.ErrorResponse](resp.Body)
	assert.Contains(t, errorBody.Error, "password is not equal to ")
}

func TestAuthService_CreateToken_EmptyUsername(t *testing.T) {

	payload := framework.CreateCustomCredentialsReqBody("", Configuration.Password)
	resp, err := AuthorizationApi.CreateToken(payload)
	assert.Equal(t, 401, resp.StatusCode)
	assert.Nil(t, err.Error)
	errorBody, _ := framework.ParseResponseBody[models.ErrorResponse](resp.Body)
	assert.Equal(t, "Invalid username or password", errorBody.Error)
}

func TestAuthService_CreateToken_EmptyPassword(t *testing.T) {

	payload := framework.CreateCustomCredentialsReqBody(Configuration.Username, "")
	resp, err := AuthorizationApi.CreateToken(payload)
	assert.Equal(t, 401, resp.StatusCode)
	assert.Nil(t, err.Error)
	errorBody, _ := framework.ParseResponseBody[models.ErrorResponse](resp.Body)
	assert.Contains(t, errorBody.Error, "password is not equal to ")
}

func TestAuthService_CreateToken_CheckTokenValidity(t *testing.T) {

	payload := framework.CreateCredentialsReqBody(Configuration.Username, Configuration.Password)
	resp, err := AuthorizationApi.CreateToken(payload)
	assert.Nil(t, err.Error)
	token, _ := framework.ParseResponseBody[server.TokenResponse](resp.Body)
	assert.True(t, len(token.Token) > 5)
	tokenisValid, _ := framework.TokenHasUsernameClaim(token.Token, Configuration.Username)
	assert.True(t, tokenisValid)

}
