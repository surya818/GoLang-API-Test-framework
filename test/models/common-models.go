package models

import (
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
