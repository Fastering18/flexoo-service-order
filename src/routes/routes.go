package routes

import (
	"os"
	"path/filepath"
	"sync"

	"flexoo-service-order/src/config"
	"flexoo-service-order/src/controllers"
	"flexoo-service-order/src/middlewares"
	"github.com/gin-gonic/gin"
)

var (
	routerInstance *gin.Engine
	routerOnce     sync.Once
)

func SetupRouter() *gin.Engine {
	routerOnce.Do(func() {
		config.ConnectDatabase()

		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(gin.Recovery())

		viewPatterns := []string{
			"src/views/*.html",
			"views/*.html",
			"../src/views/*.html",
			"/var/task/src/views/*.html",
		}
		loadedViews := false
		for _, pattern := range viewPatterns {
			if matches, err := filepath.Glob(pattern); err == nil && len(matches) > 0 {
				r.LoadHTMLGlob(pattern)
				loadedViews = true
				break
			}
		}
		if !loadedViews {
			r.SetHTMLTemplate(loadFallbackTemplates())
		}

		staticDirs := []string{
			"src/public",
			"public",
			"../src/public",
			"/var/task/src/public",
		}
		loadedStatic := false
		for _, dir := range staticDirs {
			if _, err := os.Stat(dir); err == nil {
				r.Static("/static", dir)
				loadedStatic = true
				break
			}
		}
		if !loadedStatic {
			r.GET("/static/*filepath", serveFallbackStatic)
		}

		authController := controllers.NewAuthController()
		serviceController := controllers.NewServiceController()
		orderController := controllers.NewOrderController()

		r.GET("/login", authController.ShowLogin)
		r.POST("/login", authController.Login)
		r.POST("/logout", authController.Logout)
		r.GET("/logout", authController.Logout)

		protected := r.Group("/")
		protected.Use(middlewares.AuthRequired())
		{
			protected.GET("", serviceController.Explore)
			protected.GET("services/:id/order", serviceController.ConfigureOrder)
			protected.POST("orders", orderController.PlaceOrder)
			protected.POST("orders/:id/status", orderController.UpdateStatus)
		}

		routerInstance = r
	})

	return routerInstance
}
