package engine

import (
	"github.com/backend-service/api/v1/controllers"
	"github.com/backend-service/dataservices"
	"github.com/backend-service/middleware"
	"github.com/gin-gonic/gin"
)

// BuildGinEngine creates the Gin Engine with all the middlewares, groups and routes.
func BuildGinEngine(db dataservices.BackendServiceDBInterface, version string) *gin.Engine {
	// create the default Gin engin (GIN_MODE needs to be set beforehand)
	router := gin.New()
	// attach these middlewares at root level, they will apply to every request
	router.Use(
		// recover from panics
		gin.Recovery(),
	)

	// create the /api/v1 sub-router
	v1 := router.Group("/api/v1")
	{
		v1.POST("/onboard", controllers.UserOnboard(db))
		v1.GET("/add-user/:phone", middleware.ExtractTokenMetadata(), controllers.AddUser(db))
		v1.GET("/get-all-message", middleware.ExtractTokenMetadata(), controllers.GetMessage(db))
		v1.POST("/send-message", middleware.ExtractTokenMetadata(), controllers.SendMessage(db))
		v1.GET("/get-all-user", middleware.ExtractTokenMetadata(), controllers.GetAllUser(db))
		// v1.GET("/healthcheck", xmmController.HealthCheck(db))
	}
	return router
}
