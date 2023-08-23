package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"web-ssh-server/config"
	"web-ssh-server/routes"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

func setup() *gin.Engine {
	// 1. New gin router.
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 2. Add middlewares.
	router.Use(gin.Logger())
	router.Use(Cors())

	// 3. Static resource handler.
	router.LoadHTMLGlob("**/*.html")
	router.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.Static("/static", "./static")

	// 4. Register api group.
	apiGroup := router.Group("/webssh")
	routes.SetApiRoutes(apiGroup)

	return router
}

func Run() {
	// 1. Setup router.
	router := setup()

	// 2. Create http server.
	server := &http.Server{
		Addr:    ":" + config.GlobalConfig.Port,
		Handler: router,
	}

	// 3. Startup server.
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Info("[Server] server start fail\n", err)
		}
	}()

	// 4. Wait for shutdown.
	quit := make(chan os.Signal)
	<-quit
	logrus.Info("Shutdown...")
}
