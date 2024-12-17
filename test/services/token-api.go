package service

import (
	"fmt"
	"net/http"

	"github.com/kong/candidate-take-home-exercise-sdet/internal/server"
	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"go.uber.org/zap"
)

type TokensService struct {
	Client  framework.Client
	BaseURL string
	Logger  zap.Logger
}

func NewTokensService(client framework.Client, baseUrl string) *TokensService {
	return &TokensService{
		Client:  client,
		BaseURL: baseUrl,
		Logger:  zap.Logger{},
	}
}

func (s *TokensService) CreateToken(req server.Credentials) (http.Response, framework.ApiError) {
	url := s.BaseURL + "/v1/token"
	fmt.Println("Request URL " + url)
	credentialsRequest, error := framework.StructToReader(req)
	if error != nil {
		fmt.Printf("Invalid request payload - %v", error)
	}
	resp, err := s.Client.HttpPost(url, "", credentialsRequest)

	return *resp, err

}

func (s *TokensService) FetchToken(username string, password string) string {

	var creds = server.Credentials{Username: username, Password: password}
	resp, _ := s.CreateToken(creds)
	token, _ := framework.ParseResponseBody[server.TokenResponse](resp.Body)
	return token.Token

}
