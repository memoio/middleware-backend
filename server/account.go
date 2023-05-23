package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	p.GET("/balance", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: s.NonceManager.GetNonce(),
				Error: errRes})
			return
		}
		// address := c.Query("address")
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
		// tokenString := c.GetHeader("Authorization")
		// address, err := VerifyAccessToken(tokenString)
		// if err != nil {
		// 	errRes := logs.ToAPIErrorCode(err)
		// 	c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
		// 		Nonce: s.NonceManager.GetNonce(),
		// 		Error: errRes})
		// 	return
		// }
		address := c.Query("address")

		si, err := s.Controller.GetStorageInfo(c.Request.Context(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, si)
	})
}
