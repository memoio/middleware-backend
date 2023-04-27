package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/contract"
	"github.com/memoio/backend/gateway"
	"github.com/memoio/backend/internal/gateway/ipfs"
	"github.com/memoio/backend/internal/gateway/mefs"
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
		var sr StorageResponse

		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		sr.Address = address

		api, err := mefs.NewGateway()
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		si, err := api.GetPkgSize(c.Request.Context(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		sr.StorageList = append(sr.StorageList, si)

		api = ipfs.NewGateway(s.Config.Storage.Ipfs.Host)
		si, err = api.GetPkgSize(c.Request.Context(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		sr.StorageList = append(sr.StorageList, si)

		c.JSON(http.StatusOK, sr)
	})
}
