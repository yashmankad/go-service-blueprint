package router

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"test_service/controllers"
	"test_service/repository"
)

// NewRouter initializes a new API router based on Gin
// also registers API endpoints and their handlers with the router
func NewRouter(fh *os.File, repo *repository.Repository, logger *log.Entry) (*gin.Engine, error) {
	// write API logs to the server's logfile
	// XXX: if these logs become too chatty, we may have to remove this
	// 		or write to a separate file
	gin.DefaultWriter = io.MultiWriter(fh)

	r := gin.Default()
	//gin.SetMode(gin.ReleaseMode)

	// create an instance of the controller
	ctrl := controllers.NewController(repo, logger)

	// add routes
	r.GET("/v1/ping", ctrl.Ping)

	return r, nil
}
