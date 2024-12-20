package service

import (
	"fmt"
	"net/http"

	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
)

type ServiceApi struct {
	Client    framework.Client
	BaseURL   string
	AuthToken string
}

func NewServiceApi(client framework.Client, baseUrl string, token string) *ServiceApi {
	return &ServiceApi{
		Client:    client,
		BaseURL:   baseUrl,
		AuthToken: token,
	}
}

func (s *ServiceApi) CreateService(req models.Service) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services", s.BaseURL)
	framework.Logger.Info(fmt.Sprintln("Request URL " + url))
	servicePayload, error := framework.StructToReader(req)
	if error != nil {
		framework.Logger.Info(fmt.Sprintf("Invalid request payload - %v", error))
	}
	resp, err := s.Client.HttpPost(url, s.AuthToken, servicePayload)

	return *resp, err

}

func (s *ServiceApi) UpdateService(serviceId string, req models.Service) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%v", s.BaseURL, serviceId)
	framework.Logger.Info(fmt.Sprintf("Request URL " + url))
	servicePayload, error := framework.StructToReader(req)
	if error != nil {
		framework.Logger.Info(fmt.Sprintf("Invalid request payload - %v", error))
	}
	resp, err := s.Client.HttpPatch(url, s.AuthToken, servicePayload)

	return *resp, err

}

func (s *ServiceApi) GetService(serviceId string) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%s", s.BaseURL, serviceId)
	framework.Logger.Info(fmt.Sprintln("Request URL " + url))
	resp, err := s.Client.HttpGet(url, s.AuthToken)

	return *resp, err

}

func (s *ServiceApi) DeleteService(serviceId string) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%s", s.BaseURL, serviceId)
	framework.Logger.Info(fmt.Sprintln("Request URL " + url))
	resp, err := s.Client.HttpDelete(url, s.AuthToken)

	return *resp, err

}

func (s *ServiceApi) ListServices() (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services", s.BaseURL)
	framework.Logger.Info(fmt.Sprintln("Request URL " + url))
	resp, err := s.Client.HttpGet(url, s.AuthToken)

	return *resp, err

}
