package main

import (
	"fmt"
	"kano/internal/config"
	"kano/internal/logger"
	"kano/internal/middleware"
	"kano/router"
	"strconv"

	"github.com/gin-contrib/cors"

	_ "kano/docs" // Import the docs package to generate Swagger documentation

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Kano 上传服务 API
// @version         1.0
// @description     这是 Kano 系统的通用上传服务，支持多种云存储（本地、腾讯云、阿里云）。
// @termsOfService  https://kano.com/terms/
// @contact.name   API Support
// @contact.url    https://kano.com/support
// @contact.email  support@kano.com
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
// @host      0.0.0.0:8080
// @BasePath  /kano/v1
func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize the database
	config.InitDB()

	// Create a new Gin engine
	engine := gin.New()

	logger.Std = logger.New("kano").Caller(4)
	// Use middleware
	engine.Use(gin.Recovery())
	engine.Use(logger.InitGinLogger())
	engine.Use(middleware.LoginAuth())

	engine.Use(cors.Default()) // Use CORS middleware for handling cross-origin requests

	// Initialize the router
	router.LoadRouter(engine)

	// Set the mode to release for production
	gin.SetMode(gin.ReleaseMode)

	// Start the server
	addr := config.Config.Server.Host + ":" + strconv.Itoa(config.Config.Server.Port)

	// url := ginSwagger.URL("http://47.98.202.31:9197/doc.json") // swagger.json 地址
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// logger.Std.Infof("🚀 Server starting on: %s", addr)
	fmt.Printf(`
		###########################################################################
		// +----------------------------------------------------------------------
		// | kano 上传服务
		// +----------------------------------------------------------------------
		###########################################################################
	`)
	fmt.Printf("Server is running at %s\n", addr)

	if err := engine.Run(addr); err != nil {
		panic("Failed to start server: " + err.Error())
	}

}
