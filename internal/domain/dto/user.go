package dto

import (
	"medods/internal/lib/customerror"
	"medods/internal/lib/validator"
)

type User struct {
	Email string `json:"email" validate:"email,required" example:"test@test.com"`
}

func (u *User) Validate() customerror.CustomError {
	return validator.Validate(u)
}
