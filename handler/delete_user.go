package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	if !mockUsers.delete(c.Param("id")) {
		writeError(c, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}
