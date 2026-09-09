package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename":     file.Filename,
		"size":         file.Size,
		"content_type": file.Header.Get("Content-Type"),
	})
}
