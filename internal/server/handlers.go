// Copyright © 2024 Kong Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/kong/candidate-take-home-exercise-sdet/internal/config"
	"go.uber.org/zap"
)

// Credentials represents the structure for the credentials provided by the user.
type Credentials struct {
	// Username for authentication.
	Username string `json:"username"`
	// Password for authentication.
	Password string `json:"password"`
}

// TokenResponse represents the response structure for generated JWT tokens.
type TokenResponse struct {
	// JWT Token to be used for authenticated requests.
	Token string `json:"token"`
}

// Error implements error.
func (t TokenResponse) Error() string {
	panic("unimplemented")
}

// Service represents a service in the catalog.
type Service struct {
	// Unique identifier for the service.
	ID string `json:"id"`
	// Name of the service.
	Name nullString `json:"name"`
	// Description of the service.
	Description string `json:"description"`
	// Timestamp when the service was created.
	CreatedAt time.Time `json:"created_at"`
	// Timestamp when the service was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// ServiceVersion represents a version of a specific service.
type ServiceVersion struct {
	// Unique identifier for the service version.
	ID string `json:"id"`
	// ID of the service that this version belongs to.
	ServiceID string `json:"service_id"`
	// Version information for the service.
	Version string `json:"version"`
	// Timestamp when the service version was created.
	CreatedAt time.Time `json:"created_at"`
	// Timestamp when the service version was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}

// Opts are the options used to create a new handler.
type Opts struct {
	// Config is the configuration for the application.
	Config *config.Config
	// Database is the database to use for retrieving/storing data.
	Database *sql.DB
	// Logger is the logger to use for logging.
	Logger *zap.Logger
}

// NullString is a wrapper around sql.NullString that handles JSON serialization.
type nullString struct {
	sql.NullString
}

// MarshalJSON implements the json.Marshaler interface.
// It converts a NullString to JSON null if it is not valid.
func (n nullString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	b, err := json.Marshal(n.String)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal string: %w", err)
	}
	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It converts JSON null or a string value to a NullString.
func (n *nullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		// If the JSON value is null, set the NullString to be invalid.
		n.Valid = false
		n.String = ""
		return nil
	}

	// If the JSON value is a valid string, unmarshal it and mark it as valid.
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("unable to marshal string: %w", err)
	}

	n.Valid = true
	n.String = str
	return nil
}

// Handler instance.
type Handler struct {
	jwtSecret       string
	jwtTokenTimeout time.Duration
	username        string
	password        string

	db     *sql.DB
	logger *zap.Logger
}

// NewHandler creates an instance of the handlers for the application server.
func NewHandler(opts Opts) (*Handler, error) {
	return &Handler{
		jwtSecret:       opts.Config.JWTSecret,
		jwtTokenTimeout: opts.Config.JWTTokenTimeout,
		username:        opts.Config.Username,
		password:        opts.Config.Password,

		db:     opts.Database,
		logger: opts.Logger.With(zap.String("component", "handler")),
	}, nil
}

// AuthenticateToken authenticates a request using a bearer token in the Authorization header.
func (h *Handler) AuthenticateToken(r *http.Request) error {
	// Parse auth and get the token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.logger.Warn("missing Authorization header")
		return errors.New("missing authorization header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		h.logger.Warn("invalid Authorization header format")
		return errors.New("invalid authorization header format")
	}
	tokenStr := parts[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			h.logger.Warn("unexpected signing method")
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		h.logger.Warn("invalid token", zap.Error(err))
		return errors.New("invalid token")
	}

	return nil
}

// GenerateTokenHandler handles token generation requests.
func (h *Handler) CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON payload from the request body.
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		h.logger.Error("invalid request payload", zap.Error(err))
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Validate credentials against application username and password
	if creds.Username != h.username || creds.Password != h.password {
		h.logger.Warn("invalid login attempt", zap.String("username", creds.Username))
		if creds.Password != h.password {
			http.Error(w, fmt.Sprintf(`{"error": "password is not equal to %s"}`, h.password), http.StatusUnauthorized)
		} else {
			http.Error(w, `{"error": "Invalid username or password"}`, http.StatusUnauthorized)
		}
		return
	}

	// Create a new JWT token with an expiration time.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(h.jwtTokenTimeout).Unix(),
	})

	// Sign the token using the JWT secret key.
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		h.logger.Error("failed to sign token", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Return the generated token in the response.
	response := TokenResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logger.Error("unable to encode response: %w", zap.Error(err))
	}
}

// CreateServiceHandler handles the creation of a new service in the catalog.
// It takes service information from the request and inserts it into the database.
func (h *Handler) CreateServiceHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Decode the JSON payload from the request body.
	var newService Service
	err := json.NewDecoder(r.Body).Decode(&newService)
	if err != nil {
		h.logger.Error("invalid request payload", zap.Error(err))
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Generate a new UUID v1 for the service ID.
	id, err := uuid.NewUUID()
	if err != nil {
		h.logger.Error("failed to generate UUID for new service", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	newService.ID = id.String()

	// Prepare an SQL statement to insert the new service.
	//nolint:lll
	stmt, err := h.db.Prepare("INSERT INTO services (id, name, description, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
	if err != nil {
		h.logger.Error("failed to prepare statement", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement with the service details.
	_, err = stmt.Exec(newService.ID, newService.Name, newService.Description)
	if err != nil {
		h.logger.Error("failed to insert service", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	row := h.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM services WHERE id = ?", newService.ID)

	var service Service
	err = row.Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt, &service.UpdatedAt)
	if err != nil {
		h.logger.Error("failed to fetch inserted service", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Set the response status to 201 Created and encode the new service as JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]interface{}{"item": service})
	if err != nil {
		h.logger.Error("unable to encode response", zap.Error(err))
	}
}

// ListServicesHandler lists all the available services in the catalog.
// It retrieves service information from the database and returns it in the response.
func (h *Handler) ListServicesHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Query the database to retrieve all services.
	rows, err := h.db.Query("SELECT id, name, description, created_at, updated_at FROM services")
	if err != nil {
		h.logger.Error("failed to query services", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and build a list of services.
	var services []Service
	for rows.Next() {
		var s Service
		err = rows.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			h.logger.Error("failed to scan service", zap.Error(err))
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
		services = append(services, s)
	}

	// Return the list of services in the response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"items": services})
	if err != nil {
		h.logger.Error("unable to encode response", zap.Error(err))
	}
}

// GetServiceHandler retrieves a specific service by its ID.
// It takes the service ID from the URL parameters and retrieves the service details from the database.
func (h *Handler) GetServiceHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get the service ID from the URL path variables.
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]

	// Query the database to get the service details by ID.
	var service Service
	err := h.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM services WHERE id = ?",
		serviceID).Scan(&service.ID, &service.Name, &service.Description, &service.CreatedAt, &service.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return
	} else if err != nil {
		h.logger.Error("failed to query service", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Return the service details in the response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"item": service})
	if err != nil {
		h.logger.Error("unable to encode response", zap.Error(err))
	}
}

// UpdateServiceHandler updates an existing service in the catalog.
// It takes the service ID from the URL parameters and the updated data from the request body.
func (h *Handler) UpdateServiceHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get the service ID from the URL path variables.
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]

	var updatedService Service
	updatedService.ID = serviceID

	// Decode the JSON payload from the request body.
	err := json.NewDecoder(r.Body).Decode(&updatedService)
	if err != nil {
		h.logger.Error("invalid request payload", zap.Error(err))
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	updateFields := []string{}
	values := []interface{}{}

	if updatedService.Name.Valid && strings.TrimSpace(updatedService.Name.String) != "" {
		updateFields = append(updateFields, "name = ?")
		values = append(values, updatedService.Name.String)
	}

	if strings.TrimSpace(updatedService.Description) != "" {
		updateFields = append(updateFields, "description = ?")
		values = append(values, updatedService.Description)
	}
	updateFields = append(updateFields, "updated_at = CURRENT_TIMESTAMP")
	values = append(values, serviceID)
	query := fmt.Sprintf("UPDATE services SET %s WHERE id = ?", strings.Join(updateFields, ", "))

	// Prepare an SQL statement to update the service.
	stmt, err := h.db.Prepare(query)
	if err != nil {
		h.logger.Error("failed to prepare update statement", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement with the updated service details.
	_, err = stmt.Exec(values...)
	if err != nil {
		h.logger.Error("failed to update service", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Query the database to get the service details by ID.
	var name nullString
	var description string
	err = h.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM services WHERE id = ?",
		serviceID).Scan(&updatedService.ID, &name, &description,
		&updatedService.CreatedAt, &updatedService.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, `{"error": "Service not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("failed to query service", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	if len(strings.TrimSpace(updatedService.Name.String)) == 0 {
		updatedService.Name = name
	}
	if len(strings.TrimSpace(updatedService.Description)) == 0 {
		updatedService.Description = description
	}

	// Return a 200 OK response indicating the service was successfully updated.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"item": updatedService})
	if err != nil {
		h.logger.Error("unable to encode response", zap.Error(err))
	}
}

// DeleteServiceHandler deletes a specific service from the catalog.
// It takes the service ID from the URL parameters and removes the service from the database.
func (h *Handler) DeleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get the service ID from the URL path variables.
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]

	// Prepare an SQL statement to delete the service by ID.
	stmt, err := h.db.Prepare("DELETE FROM services WHERE id = ?")
	if err != nil {
		h.logger.Error("failed to prepare delete statement", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the service.
	_, err = stmt.Exec(serviceID)
	if err != nil {
		h.logger.Error("failed to delete service", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Return a 204 No Content response indicating the service was successfully deleted.
	w.WriteHeader(http.StatusNoContent)
}

// CreateServiceVersionHandler handles the creation of a new version for a specific service.
// It takes the service ID from the URL and the version data from the request body, and inserts it into the database.
func (h *Handler) CreateServiceVersionHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get the service ID from the URL path variables.
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]

	var newVersion ServiceVersion

	// Decode the JSON payload from the request body.
	err := json.NewDecoder(r.Body).Decode(&newVersion)
	if err != nil {
		h.logger.Error("invalid request payload", zap.Error(err))
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Generate a new UUID version ID.
	id, err := uuid.NewUUID()
	if err != nil {
		h.logger.Error("failed to generate UUID for new service version", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	newVersion.ID = id.String()
	newVersion.ServiceID = serviceID
	version := newVersion.Version

	if len(version) > 16 {
		h.logger.Error("invalid request payload", zap.Error(err))
		http.Error(w, `{"error": "Version cannot be longer than 16 characters "}`, http.StatusBadRequest)
		return
	}

	// Prepare an SQL statement to insert the new version.
	//nolint:lll
	stmt, err := h.db.Prepare("INSERT INTO service_versions (id, service_id, version, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
	if err != nil {
		h.logger.Error("failed to prepare statement", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement with the version details.
	_, err = stmt.Exec(newVersion.ID, newVersion.ServiceID, newVersion.Version)
	if err != nil {
		h.logger.Error("failed to insert service version", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Set the response status to 201 Created and encode the new version as JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]interface{}{"item": newVersion})
	if err != nil {
		h.logger.Error("unable to encode response", zap.Error(err))
	}
}

// ListServiceVersionsHandler lists all the versions for a specific service.
// It retrieves version information from the database for the given service ID.
func (h *Handler) ListServiceVersionsHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get the service ID from the URL path variables.
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]

	// Query the database to retrieve all versions for the given service.
	//nolint:lll
	rows, err := h.db.Query("SELECT id, service_id, version, created_at, updated_at FROM service_versions WHERE service_id = ?", serviceID)
	if err != nil {
		h.logger.Error("failed to query service versions", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and build a list of versions.
	var versions []ServiceVersion
	for rows.Next() {
		var v ServiceVersion
		err = rows.Scan(&v.ID, &v.ServiceID, &v.Version, &v.CreatedAt, &v.UpdatedAt)
		if err != nil {
			h.logger.Error("failed to scan service version", zap.Error(err))
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}
		versions = append(versions, v)
	}

	// Return the list of versions in the response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"items": versions})
	if err != nil {
		h.logger.Error("unable to encode response", zap.Error(err))
	}
}

// GetServiceVersionHandler retrieves a specific version for a given service by its ID.
// It takes the service ID and version ID from the URL and retrieves the version details from the database.
func (h *Handler) GetServiceVersionHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get the service ID and version ID from the URL path variables.
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]
	versionID := vars["versionId"]

	// Query the database to get the version details by ID.
	var version ServiceVersion
	//nolint:lll
	err := h.db.QueryRow("SELECT id, service_id, version, created_at, updated_at FROM service_versions WHERE id = ? AND service_id = ?",
		versionID, serviceID).Scan(&version.ID, &version.ServiceID, &version.Version, &version.CreatedAt, &version.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, `{"error": "Service version not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("failed to query service version", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Return the version details in the response.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"item": version})
	if err != nil {
		h.logger.Error("unable to encode response", zap.Error(err))
	}
}

// UpdateServiceVersionHandler updates an existing version for a specific service.
// It takes the service ID and version ID from the URL and the updated data from the request body.
func (h *Handler) UpdateServiceVersionHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get the service ID and version ID from the URL path variables.
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]
	versionID := vars["versionId"]

	var updatedVersion ServiceVersion
	updatedVersion.ID = versionID
	updatedVersion.ServiceID = serviceID

	// Decode the JSON payload from the request body.
	err := json.NewDecoder(r.Body).Decode(&updatedVersion)
	if err != nil {
		h.logger.Error("invalid request payload", zap.Error(err))
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	version := updatedVersion.Version
	if len(version) > 16 {
		h.logger.Error("Invalid request payload")
		http.Error(w, `{"error": "Version cannot be longer than 16 characters "}`, http.StatusBadRequest)
		return
	}
	// Prepare an SQL statement to update the service version.
	//nolint:lll
	stmt, err := h.db.Prepare("UPDATE service_versions SET ID = ?, version = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND service_id = ?")
	if err != nil {
		h.logger.Error("failed to prepare update statement", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Generate a new UUID version ID.
	id, err := uuid.NewUUID()
	if err != nil {
		h.logger.Error("failed to generate UUID for new service version", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Execute the SQL statement with the updated version details.
	_, err = stmt.Exec(id.String(), updatedVersion.Version, versionID, serviceID)
	if err != nil {
		h.logger.Error("failed to update service version", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Query the database to get the version details by ID.
	err = h.db.QueryRow("SELECT created_at, updated_at FROM service_versions WHERE id = ? AND service_id = ?",
		id.String(), serviceID).Scan(&updatedVersion.CreatedAt, &updatedVersion.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, `{"error": "Service version not found"}`, http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error("failed to query service version", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Return a 200 OK response indicating the version was successfully updated.
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"item": updatedVersion})
	if err != nil {
		h.logger.Error("unable to encode response", zap.Error(err))
	}
}

// DeleteServiceVersionHandler deletes a specific version for a given service.
// It takes the service ID and version ID from the URL and removes the version from the database.
func (h *Handler) DeleteServiceVersionHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.AuthenticateToken(r); err != nil {
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Get the service ID and version ID from the URL path variables.
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]
	versionID := vars["versionId"]

	// Prepare an SQL statement to delete the service version by ID.
	stmt, err := h.db.Prepare("DELETE FROM service_versions WHERE id = ? AND service_id = ?")
	if err != nil {
		h.logger.Error("failed to prepare delete statement", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the version.
	_, err = stmt.Exec(versionID, serviceID)
	if err != nil {
		h.logger.Error("failed to delete service version", zap.Error(err))
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// Return a 204 No Content response indicating the version was successfully deleted.
	w.WriteHeader(http.StatusNoContent)
}
