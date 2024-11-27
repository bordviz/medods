package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessToken struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

type RefreshToken struct {
	jwt.RegisteredClaims
	UserID    uuid.UUID `json:"user_id"`
	IpAddress string    `json:"ip_address"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
