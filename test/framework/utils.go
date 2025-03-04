package framework

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kong/candidate-take-home-exercise-sdet/internal/config"
	"github.com/kong/candidate-take-home-exercise-sdet/internal/server"
	"github.com/kong/candidate-take-home-exercise-sdet/test/models"
	"gopkg.in/yaml.v3"
)

func ResponseBodyToString(response http.Response) string {

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error reading response body: %v\n", err))
		return ""
	}

	// Convert the byte slice to a string
	bodyString := string(bodyBytes)
	return bodyString
}

// This method is used to print request payloads in HttpPost method
func ReaderToString(r io.Reader) (string, error) {
	// Read all contents from the io.Reader
	body, err := io.ReadAll(r) // Use io.ReadAll in Go 1.16+ instead of ioutil.ReadAll
	if err != nil {
		return "", err
	}
	// Return as a string
	return string(body), nil
}

// This method is used to parse the Response body to desired type
func ParseResponseBody[T any](body io.ReadCloser) (T, error) {
	var result T
	defer body.Close()
	err := json.NewDecoder(body).Decode(&result)

	if err != nil {
		return result, fmt.Errorf("failed to parse response body: %w", err)
	}
	return result, nil
}

func StructToReader[T any](s T) (io.Reader, error) {
	// Serialize the struct to JSON
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	// Wrap the JSON data as an io.Reader
	return bytes.NewReader(data), nil
}

func GetConfiguration() config.Config {
	configFile, err := os.Open("../../config.yml")
	if err != nil {
		Logger.Info(fmt.Sprintf("Error opening config file: %v\n", err))
		return config.Config{}
	}
	defer configFile.Close()

	// Unmarshal the YAML data into the Config struct
	var configuration config.Config
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&configuration)
	if err != nil {
		Logger.Info(fmt.Sprintf("Error decoding YAML: %v\n", err))
		return config.Config{}
	}
	return configuration
}

func isJWT(token string) bool {
	// Split the token into parts by the dot separator
	parts := strings.Split(token, ".")

	// JWT should have exactly 3 parts (Header, Payload, Signature)
	if len(parts) != 3 {
		return false
	}

	// Check if both the Header and Payload parts are valid Base64
	// Base64 decoding of JWT parts should not return errors
	_, errHeader := base64.RawURLEncoding.DecodeString(parts[0])
	_, errPayload := base64.RawURLEncoding.DecodeString(parts[1])

	if errHeader != nil || errPayload != nil {
		return false
	}

	// If both Base64 decodings are successful, it's likely a valid JWT format
	return true
}

func TokenHasUsernameClaim(tokenString string, expectedUsername string) (bool, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &models.KongJWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		// Return the secret key used to sign the token (replace with your actual secret key)
		return []byte(GetConfiguration().JWTSecret), nil
	})

	// Handle errors during parsing
	if err != nil {
		return false, fmt.Errorf("failed to parse token: %v", err)
	}

	// Check if the token is valid and if the claims match
	if claims, ok := token.Claims.(*models.KongJWTClaim); ok && token.Valid {
		// Check if the username claim matches the expected username
		if claims.Username == expectedUsername {
			return true, nil
		}
	}

	// Return false if the token is invalid or the username does not match
	return false, nil
}

func GetRandomNumber() string {
	rand.Seed(time.Now().UnixNano())         // Seed the random number generator
	randomNumber := rand.Intn(90000) + 10000 // Ensure the number is 5 digits (10000–99999)
	return fmt.Sprintf("%05d", randomNumber) // Format as a 5-digit string
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func GetRandomName(prefix string) string {
	suffix := GetRandomNumber()
	return prefix + "-" + suffix
}

func CreateCredentialsReqBody(username string, password string) server.Credentials {
	creds := server.Credentials{Username: username, Password: password}
	return creds
}

func CreateCustomCredentialsReqBody(username string, password string) server.Credentials {
	creds := server.Credentials{Username: username, Password: password}
	return creds
}

func CreateServicePayload(id string, name string, description string) models.Service {

	if description == "" {
		description = "newly created service"
	}

	service := models.Service{ID: id, Name: name, Description: description}
	return service
}

func CreateServiceVersionPayload(serviceId string, id string, version string) models.ServiceVersion {

	num := GetRandomNumber()
	if version == "" {
		version = "v" + num
	}

	if id == "" {
		id = "id-" + num
	}

	serviceVersion := models.ServiceVersion{ServiceID: serviceId, ID: id, Version: version}
	return serviceVersion
}
