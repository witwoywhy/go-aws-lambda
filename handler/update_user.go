package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context) {
	input, err := bindUser(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	user, ok := mockUsers.update(c.Param("id"), input)
	if !ok {
		writeError(c, http.StatusNotFound, "user not found")
		return
	}
	c.JSON(http.StatusOK, user)
}
