package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is the standard API error envelope.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

// OK sends a 200 success response.
func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 created response.
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a 400 error response.
func BadRequest(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

// Unauthorized sends a 401 error response.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// Forbidden sends a 403 error response.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// NotFound sends a 404 error response.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// InternalServerError sends a 500 error response.
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Message: message,
	})
}

// UnprocessableEntity sends a 422 error response.
func UnprocessableEntity(c *gin.Context, message string, errors interface{}) {
	c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
