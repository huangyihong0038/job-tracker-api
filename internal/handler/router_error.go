package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context) {
	writeError(
		c,
		http.StatusNotFound,
		"NOT_FOUND",
		"route not found",
	)
}

func MethodNotAllowed(c *gin.Context) {
	writeError(
		c,
		http.StatusMethodNotAllowed,
		"METHOD_NOT_ALLOWED",
		"method not allowed",
	)
}
