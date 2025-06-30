package router

import (
	v1 "kano/internal/handler/v1"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoadRouter initializes the upload router for handling file uploads and related operations.
func LoadRouter(engine *gin.Engine) {

	engine.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Welcome to Kano service")
	})

	// 加载 HTML 模板文件
	engine.LoadHTMLGlob("templates/*")

	// 定义一个简单的 GET 路由，渲染 HTML 模板
	engine.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":   "Welcome to Gin",
			"Message": "This is a simple HTML page rendered by Gin.",
		})
	})

	// Initialize the upload router
	uploadRouter := engine.Group("/kano/v1")
	{
		uploadRouter.GET("/upload/token", v1.GetUploadToken)
		uploadRouter.POST("/upload/record", v1.SaveUploadRecord)
	}
}
