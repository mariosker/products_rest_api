package utils

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse defines a standard error format.
// @Description ErrorResponse defines the standard format for error responses.
// @Success 400 {object} ErrorResponse
// @Success 500 {object} ErrorResponse
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// SendErrorResponse sends a standardized JSON error response.
// @Param status query int true "HTTP status code"
// @Param message query string true "Error message"
// @Success 200 {object} ErrorResponse
func SendErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Status:  status,
		Message: message,
	})
}
