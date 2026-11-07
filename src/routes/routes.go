package routes

import (
	"log"
	"os"
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

		r.SetHTMLTemplate(loadFallbackTemplates())
		log.Println("[flexoo] loaded embedded HTML templates")

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
				log.Printf("[flexoo] serving static files from: %s", dir)
				loadedStatic = true
				break
			}
		}
		if !loadedStatic {
			r.GET("/static/*filepath", serveFallbackStatic)
			log.Println("[flexoo] using embedded fallback static assets")
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
		log.Println("[flexoo] router initialized successfully")
	})

	return routerInstance
}
