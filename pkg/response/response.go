package response

import (
	"github.com/gin-gonic/gin"
)

// Response adalah format standar envelope response JSON untuk seluruh endpoint API
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

// PaginationMeta adalah informasi metadata penomoran halaman untuk data list
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// Success mengirimkan response sukses standar (200, 201, dll)
func Success(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta mengirimkan response sukses dengan metadata pagination
func SuccessWithMeta(c *gin.Context, statusCode int, message string, data any, meta any) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Error mengirimkan response gagal/error standar (400, 404, 422, 500, dll)
func Error(c *gin.Context, statusCode int, message string, errors any) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
