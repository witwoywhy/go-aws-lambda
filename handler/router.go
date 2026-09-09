package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/witwoywhy/go-cores/gins"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gins.Log(nil))

	router.GET("/users", SelectUsers)
	router.GET("/users/:id", SelectUser)
	router.POST("/users", CreateUser)
	router.PUT("/users/:id", UpdateUser)
	router.DELETE("/users/:id", DeleteUser)
	router.POST("/files", UploadFile)

	return router
}
