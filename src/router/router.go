package router

import (
	"github.com/gin-gonic/gin"

	v1api "test_service/api/v1"
)

// initializes a new API router based on Gin
// also registers API endpoints and their handlers with the router
func NewRouter() (*gin.Engine, error) {
	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)

	// add routes
	r.GET("/v1/ping", v1api.Ping)

	return r, nil
}
