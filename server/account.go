package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
)

type StorageResponse struct {
	Address     string
	StorageList []storage.StorageInfo
}

func (s Server) accountRegistRoutes(r *gin.RouterGroup) {
	s.addGetBalanceRoutes(r)
	s.addGetStorageRoutes(r)
	s.addBuyPkgRoutes(r)
	s.addGetPkgListRoutes(r)
	s.addGetBuyPkgRoutes(r)
}

func (s Server) addGetBalanceRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/balance", auth.VerifyIdentityHandler, func(c *gin.Context) {
		address := c.GetString("address")
		balance, err := s.Controller.GetBalance(c.Request.Context(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		c.JSON(http.StatusOK, BalanceResponse{Address: address, Balance: balance.String()})
	})
}

func (s Server) addGetStorageRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/storageinfo", func(c *gin.Context) {
		address := c.GetString("address")

		si, err := s.Controller.GetStorageInfo(c.Request.Context(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, si)
	})
}
