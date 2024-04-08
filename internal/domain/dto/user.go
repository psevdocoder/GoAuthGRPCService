package dto

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
)

type LoginInput struct {
	Username string `validate:"required,min=4,max=32"`
	Password string `validate:"required,min=8,max=48,alphanum"`
	AppID    uint32 `validate:"required,numeric,gt=0"`
}

func (i LoginInput) Validate() error {
	validate := validator.New(validator.WithPrivateFieldValidation())
	return validate.Struct(i)
}

type RegisterInput struct {
	Username string `validate:"required,min=4,max=32"`
	Password string `validate:"required,alphanum,min=8"`
}

func (v *RegisterInput) Validate() error {
	validate := validator.New(validator.WithPrivateFieldValidation())
	return validate.Struct(v)
}

type GetRoleInput struct {
	Username string
}
