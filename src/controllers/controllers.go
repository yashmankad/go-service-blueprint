package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"test_service/models"
	"test_service/repository"
)

// Controller is a wrapper object for all API handlers
type Controller struct {
	// Repository object (includes conn object to the db/repo)
	Repository *repository.Repository

	// Logger object
	Logger *log.Entry
}

// NewController will create a new controller object
func NewController(repo *repository.Repository, logger *log.Entry) Controller {
	return Controller{
		Repository: repo,
		Logger:     logger,
	}
}

// Ping API endpoint handler to verify server's REST server availability
func (ctrl *Controller) Ping(c *gin.Context) {
	var response models.PingResponse
	response.Message = "pong"
	c.JSON(http.StatusOK, &response)
}
