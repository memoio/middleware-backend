package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	global "github.com/memoio/backend/global/db"
)

func DBHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.InitDB()
		c.JSON(http.StatusOK, "init db")
	}
}
