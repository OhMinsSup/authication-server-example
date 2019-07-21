package lib

import (
	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}
