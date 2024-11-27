package dto

import (
	"fmt"
	"medods/internal/lib/validator"
)

type User struct {
	Email string `json:"email" validate:"email,required"`
}

func (u *User) Validate() error {
	if err := validator.Validate(u); err != "" {
		return fmt.Errorf("validation error: %s", err)
	}
	return nil
}
