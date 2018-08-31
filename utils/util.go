package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckError(c *gin.Context, err error){
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

