package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SelectUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"users": mockUsers.list()})
}

func SelectUser(c *gin.Context) {
	user, ok := mockUsers.get(c.Param("id"))
	if !ok {
		writeError(c, http.StatusNotFound, "user not found")
		return
	}
	c.JSON(http.StatusOK, user)
}
