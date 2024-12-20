package service

import (
	"fmt"
	"net/http"

	"github.com/kong/candidate-take-home-exercise-sdet/test/framework"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
	"go.uber.org/zap"
)

type ServiceVersionApi struct {
	Client    framework.Client
	BaseURL   string
	Logger    zap.Logger
	AuthToken string
}

func NewServiceVersionApi(client framework.Client, baseUrl string, token string) *ServiceVersionApi {
	return &ServiceVersionApi{
		Client:    client,
		BaseURL:   baseUrl,
		Logger:    zap.Logger{},
		AuthToken: token,
	}
}

func (s *ServiceVersionApi) CreateServiceVersion(serviceId string, req models.ServiceVersion) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%v/versions", s.BaseURL, serviceId)
	framework.Logger.Info(fmt.Sprintf("Request URL " + url))
	serviceVersionPayload, error := framework.StructToReader(req)
	if error != nil {
		framework.Logger.Error(fmt.Sprintf("Invalid request payload - %v", error))
	}
	resp, err := s.Client.HttpPost(url, s.AuthToken, serviceVersionPayload)

	return *resp, err

}

func (s *ServiceVersionApi) UpdateServiceVersion(serviceId string, versionId string, req models.ServiceVersion) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%v/versions/%v", s.BaseURL, serviceId, versionId)
	framework.Logger.Info(fmt.Sprintf("Request URL " + url))
	serviceVersionPayload, error := framework.StructToReader(req)
	if error != nil {
		framework.Logger.Error(fmt.Sprintf("Invalid request payload - %v", error))
	}
	resp, err := s.Client.HttpPatch(url, s.AuthToken, serviceVersionPayload)

	return *resp, err

}

func (s *ServiceVersionApi) GetServiceVersion(serviceId string, versionId string) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%v/versions/%v", s.BaseURL, serviceId, versionId)
	framework.Logger.Info(fmt.Sprintf("Request URL " + url))
	resp, err := s.Client.HttpGet(url, s.AuthToken)

	return *resp, err

}

func (s *ServiceVersionApi) DeleteServiceVersion(serviceId string, versionId string) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%v/versions/%v", s.BaseURL, serviceId, versionId)
	framework.Logger.Info(fmt.Sprintf("Request URL " + url))
	resp, err := s.Client.HttpDelete(url, s.AuthToken)

	return *resp, err

}

func (s *ServiceVersionApi) ListServiceVersions(serviceId string) (http.Response, framework.ApiError) {
	url := fmt.Sprintf("%s/v1/services/%v/versions", s.BaseURL, serviceId)
	framework.Logger.Info(fmt.Sprintf("Request URL " + url))
	resp, err := s.Client.HttpGet(url, s.AuthToken)

	return *resp, err

}
