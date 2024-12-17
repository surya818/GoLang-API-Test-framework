package models

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type KongJWTClaim struct {
	Username string `json:"username"`
	Expiry   int    `json:"exp"`
	jwt.RegisteredClaims
}

// Service represents a service in the catalog.
type Service struct {
	// Unique identifier for the service.
	ID string `json:"id"`
	// Name of the service.
	Name string `json:"name"`
	// Description of the service.
	Description string `json:"description"`
	// Timestamp when the service was created.
	CreatedAt time.Time `json:"created_at"`
	// Timestamp when the service was last updated.
	UpdatedAt time.Time `json:"updated_at"`
}
type NullString struct {
	sql.NullString
}

type ServiceResponse struct {
	Item Service `json:"item"`
}

type ListServices struct {
	Items []Service `json:"items"`
}
