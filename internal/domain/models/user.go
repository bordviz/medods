package models

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refresh_token"`
}

type CreateUser struct {
	Detail string `json:"detail" example:"new user was successfully created"`
	ID     string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}
