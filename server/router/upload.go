package router

import (
	v1 "kano/internal/handler/v1"

	"github.com/gin-gonic/gin"
)

// LoadRouter initializes the upload router for handling file uploads and related operations.
func LoadRouter(engine *gin.Engine) {
	// Initialize the upload router
	uploadRouter := engine.Group("/knano/v1")
	{
		uploadRouter.GET("/upload/token", v1.GetUploadToken)
		uploadRouter.POST("/upload/record", v1.SaveUploadRecord)
	}
}
