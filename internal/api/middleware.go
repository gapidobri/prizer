package api

import (
	"fmt"
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	er "github.com/gapidobri/prizer/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			return field.Tag.Get("json")
		})
	}
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}

	err := c.Errors[0].Err

	var apiErr er.ApiError
	if errors.As(err, &apiErr) {
		c.JSON(apiErr.StatusCode(), api.ErrorResponse{
			Error: apiErr.Message(),
			Code:  apiErr.Code(),
		})
		return
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		err := validationErrors[0]

		errString := fmt.Sprintf("Field '%s' ", err.Field())
		switch err.Tag() {
		case "required":
			errString += "is missing"
		case "email":
			errString += "is not a valid email address"
		case "uuid":
			errString += "is not a valid UUID"
		default:
			errString += fmt.Sprintf("failed on '%s'", err.Tag())
		}

		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: errString,
			Code:  "validation_error",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, api.ErrorResponse{
		Error: "An unknown error occurred",
		Code:  "internal_server_error",
	})
}
