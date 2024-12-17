package framework

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Client interface {
	HttpGet(url string, token string) (*http.Response, ApiError)
	HttpPost(url string, token string, payload io.Reader) (*http.Response, ApiError)
	HttpDelete(url string, token string) (*http.Response, ApiError)
	HttpPatch(url string, token string, payload io.Reader) (*http.Response, ApiError)
}
type ApiError struct {
	Error    error
	Response *http.Response
}
type Todo struct {
	ID        int    `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
	UserID    int    `json:"userId"`
}

type HttpClient struct {
	baseUrl    string
	HttpClient *http.Client
	Logger     *zap.Logger
}

func NewHttpClient(url string) *HttpClient {

	return &HttpClient{
		baseUrl:    url,
		HttpClient: &http.Client{},
	}
}

func (httpClient *HttpClient) HttpGet(path string, token string) (*http.Response, ApiError) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		errMsg := fmt.Sprintf("Error building the request %v with reason %v ", path, err.Error())
		fmt.Println(errMsg)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	fmt.Printf("\nInvoking %v %v \n", req.Method, req.URL)
	resp, err := httpClient.HttpClient.Do(req)
	apierror := ApiError{err, resp}
	if err != nil {
		errMsg := fmt.Sprintf("Error invoking the GET %v with reason %v ", path, err.Error())
		fmt.Println(errMsg)
	}
	logResponse(resp)
	return resp, apierror
}

func (httpClient *HttpClient) HttpDelete(path string, token string) (*http.Response, ApiError) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
	if err != nil {
		errMsg := fmt.Sprintf("Error building the request %v with reason %v ", path, err.Error())
		fmt.Println(errMsg)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	fmt.Printf("\nInvoking %v %v \n", req.Method, req.URL)
	resp, err := httpClient.HttpClient.Do(req)
	apierror := ApiError{err, resp}
	if err != nil {
		errMsg := fmt.Sprintf("Error invoking the GET %v with reason %v ", path, err.Error())
		fmt.Println(errMsg)
	}
	logResponse(resp)
	return resp, apierror
}

func (httpClient *HttpClient) HttpPost(path string, token string, payload io.Reader) (*http.Response, ApiError) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, payload)
	if err != nil {
		errMsg := fmt.Sprintf("Error building the request %v with reason %v ", path, err.Error())
		fmt.Println(errMsg)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rc, _ := req.GetBody()
	payloadvalue, _ := ReaderToString(rc)
	fmt.Printf("\nInvoking %v %v \nPayload: %v", req.Method, req.URL, payloadvalue)
	resp, err := httpClient.HttpClient.Do(req)
	apierror := ApiError{err, resp}
	if err != nil {
		errMsg := fmt.Sprintf("\nError invoking the POST %v with reason %v \n", path, err.Error())
		fmt.Println(errMsg)
	} else {
		// Create a buffer to hold a copy of the response body
		// Read the response body into bodyBuffer
		// Save the body into the buffer
		// Reassign the body to allow further reading
		// Print the body content for debugging/logging
		logResponse(resp)
	}

	return resp, apierror
}

func (httpClient *HttpClient) HttpPatch(path string, token string, payload io.Reader) (*http.Response, ApiError) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, path, payload)
	if err != nil {
		errMsg := fmt.Sprintf("Error building the request %v with reason %v ", path, err.Error())
		fmt.Println(errMsg)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rc, _ := req.GetBody()
	payloadvalue, _ := ReaderToString(rc)
	fmt.Printf("\nInvoking %v %v \nPayload: %v", req.Method, req.URL, payloadvalue)
	resp, err := httpClient.HttpClient.Do(req)
	apierror := ApiError{err, resp}
	if err != nil {
		errMsg := fmt.Sprintf("\nError invoking the POST %v with reason %v \n", path, err.Error())
		fmt.Println(errMsg)
	} else {
		// Create a buffer to hold a copy of the response body
		// Read the response body into bodyBuffer
		// Save the body into the buffer
		// Reassign the body to allow further reading
		// Print the body content for debugging/logging
		logResponse(resp)
	}

	return resp, apierror
}

func logResponse(resp *http.Response) {
	fmt.Printf("\nStatus Code: %v \n", resp.StatusCode)

	var bodyBuffer bytes.Buffer

	bodyContent, readerr := io.ReadAll(resp.Body)
	if readerr != nil {
		fmt.Printf("Error performing operation: %v", readerr)
	}

	bodyBuffer.Write(bodyContent)

	resp.Body = io.NopCloser(&bodyBuffer)

	fmt.Println("Response Body \n", string(bodyContent))
}
