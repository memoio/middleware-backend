package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/gateway/ipfs"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			c.Abort()
			return
		}
	}
}

func LoadMefsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ui := api.USerInfo{
			Api:   config.Cfg.Storage.Mefs.Api,
			Token: config.Cfg.Storage.Mefs.Token,
		}
		store, err := mefs.NewGatewayWith(ui)
		if err != nil {
			logger.Error("init mefs error:", err)
		}
		c.Set("store", store)
	}
}

func LoadIpfsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		store, err := ipfs.NewGateway()
		if err != nil {
			logger.Error("init ipfs error:", err)
		}
		c.Set("store", store)
	}
}
