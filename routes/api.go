package routes

import (
	"github.com/gin-gonic/gin"
	"web-ssh-server/service"
	"web-ssh-server/websocket"
)

func SetApiRoutes(router *gin.RouterGroup) {
	router.GET("/ws/:sessionId", websocket.WebSshWsHandler)
	router.POST("/v1/auth/login", service.LoginServiceHandler)
	router.POST("/v1/download", service.DownloadFileHandler)
	router.POST("/v1/listFiles", service.ListFilesHandler)
	router.POST("/v1/deleteFile", service.DeleteFileHandler)
	router.POST("/v1/upload", service.UploadFileHandler)
	router.POST("/v1/getRootPath", service.GetDefaultRootPathHandler)
	router.POST("/v1/getParentPath", service.GetParentDirHandler)
}
