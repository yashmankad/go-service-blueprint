package router

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"

	"test_service/v1api"
)

// NewRouter initializes a new API router based on Gin
// also registers API endpoints and their handlers with the router
func NewRouter(fh *os.File) (*gin.Engine, error) {
	// write API logs to the server's logfile
	// XXX: if these logs become too chatty, we may have to remove this
	// 		or write to a separate file (like tomcat's api logs)
	gin.DefaultWriter = io.MultiWriter(fh)

	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)

	// add routes
	r.GET("/v1/ping", v1api.Ping)

	return r, nil
}
