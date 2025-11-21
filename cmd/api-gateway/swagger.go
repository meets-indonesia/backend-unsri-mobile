package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title UNSRI Backend API
// @version 1.0
// @description Backend API untuk aplikasi mobile UNSRI dengan arsitektur microservices
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@unsri.ac.id

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// setupSwagger initializes Swagger documentation
// Note: To generate swagger docs, run: swag init -g cmd/api-gateway/main.go
// This will create docs/ folder with swagger.json and swagger.yaml
// After generating, uncomment the docs import and SwaggerInfo setup below
func setupSwagger(router *gin.Engine) {
	// Uncomment after running 'swag init':
	// import "unsri-backend/docs"
	// docs.SwaggerInfo.Title = "UNSRI Backend API"
	// docs.SwaggerInfo.Description = "Backend API untuk aplikasi mobile UNSRI dengan arsitektur microservices"
	// docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8080"
	// docs.SwaggerInfo.BasePath = "/api/v1"
	// docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Swagger UI will be available at /swagger/index.html
	// Make sure to run 'swag init' first to generate docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
