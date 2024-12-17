package service

import (
	"fmt"
	"net/http"

	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
	"go.uber.org/zap"
)

type ServiceApi struct {
	Client    framework.Client
	BaseURL   string
	Logger    zap.Logger
	AuthToken string
}

func NewServiceApi(client framework.Client, baseUrl string, token string) *ServiceApi {
	return &ServiceApi{
		Client:  client,
		BaseURL: baseUrl,
		Logger:  zap.Logger{},
	}
}

func (s *ServiceApi) CreateService(req models.Service) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services", s.BaseURL)
	fmt.Println("Request URL " + url)
	servicePayload, error := framework.StructToReader(req)
	if error != nil {
		fmt.Printf("Invalid request payload - %v", error)
	}
	resp, err := s.Client.HttpPost(url, s.AuthToken, servicePayload)

	return *resp, err

}

func (s *ServiceApi) UpdateService(serviceId string, req models.Service) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%v", s.BaseURL, serviceId)
	fmt.Println("Request URL " + url)
	servicePayload, error := framework.StructToReader(req)
	if error != nil {
		fmt.Printf("Invalid request payload - %v", error)
	}
	resp, err := s.Client.HttpPatch(url, s.AuthToken, servicePayload)

	return *resp, err

}

func (s *ServiceApi) GetService(serviceId string) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%s", s.BaseURL, serviceId)
	fmt.Println("Request URL " + url)
	resp, err := s.Client.HttpGet(url, s.AuthToken)

	return *resp, err

}

func (s *ServiceApi) DeleteService(serviceId string) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%s", s.BaseURL, serviceId)
	fmt.Println("Request URL " + url)
	resp, err := s.Client.HttpDelete(url, s.AuthToken)

	return *resp, err

}

func (s *ServiceApi) ListServices() (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services", s.BaseURL)
	fmt.Println("Request URL " + url)
	resp, err := s.Client.HttpGet(url, s.AuthToken)

	return *resp, err

}
