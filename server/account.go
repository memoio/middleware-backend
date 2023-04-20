package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/contract"
	"github.com/memoio/backend/gateway"
	"github.com/memoio/backend/internal/storage"
)

func (s Server) addGetBalanceRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/balance", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: s.NonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		balance := contract.BalanceOf(c.Request.Context(), address)
		c.JSON(http.StatusOK, BalanceResponse{Address: address, Balance: balance.String()})
	})
}

func (s Server) addGetStorageRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/getstorage", func(c *gin.Context) {
		stype := c.Query("stype")

		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: s.NonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		si, err := s.Gateway.GetPkgSize(c.Request.Context(), storage.ToStorageType(stype), address)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, apiErr)
		}
		c.JSON(http.StatusOK, si)
	})
}
