package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DBHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "init db")
	}
}
