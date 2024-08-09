package api

import (
	"github.com/gapidobri/prizer/internal/pkg/models/api"
	er "github.com/gapidobri/prizer/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

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
		})
		return
	}

	c.JSON(http.StatusInternalServerError, api.ErrorResponse{
		Error: "An unknown error occurred",
	})
}
