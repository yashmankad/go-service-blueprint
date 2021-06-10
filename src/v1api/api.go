// v1 API package defines API endpoints and their behavior

package v1api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingResponse is the server response for the ping API endpoint
type PingResponse struct {
	Message string `json:"message"`
}

// Ping API endpoint handler to verify server's REST capabilities
func Ping(c *gin.Context) {
	var response PingResponse
	response.Message = "pong"
	c.JSON(http.StatusOK, &response)
}
