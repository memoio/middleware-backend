package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/logs"
)

func (s Server) accountRegistRoutes(r *gin.RouterGroup) {
	s.addGetBalanceRoutes(r)
	s.addGetStorageRoutes(r)
}

func (s Server) addGetBalanceRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/balance", auth.VerifyIdentityHandler, func(c *gin.Context) {
		address := c.GetString("address")
		chain := c.GetInt("chainid")
		balance, err := s.Controller.GetBalance(c.Request.Context(), chain, address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		c.JSON(http.StatusOK, gin.H{"Address": address, "Balance": balance.String()})
	})
}

func (s Server) addGetStorageRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/storageinfo", auth.VerifyIdentityHandler, func(c *gin.Context) {
		address := c.GetString("address")
		chain := c.GetInt("chainid")

		si, err := s.Controller.GetStorageInfo(c.Request.Context(), chain, address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, si)
	})
}
