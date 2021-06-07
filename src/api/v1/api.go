// v1 API package defines API endpoints and their behavior

package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingResponse struct {
	Message string `json:"message"`
}

// test API to verify http router and REST capabilities
func Ping(c *gin.Context) {
	var response PingResponse
	response.Message = "pong"
	c.JSON(http.StatusOK, &response)
}
