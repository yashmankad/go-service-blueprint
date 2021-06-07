// v1 API package defines API endpoints and their behavior

package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// test API to verify http router and REST capabilities
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
