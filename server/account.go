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
	// s.addGetStorageRoutes(r)
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

// func (s Server) addGetStorageRoutes(r *gin.RouterGroup) {
// 	p := r.Group("/")
// 	p.GET("/getstorage", func(c *gin.Context) {
// 		var sr StorageResponse

// 		tokenString := c.GetHeader("Authorization")
// 		address, err := VerifyAccessToken(tokenString)
// 		if err != nil {
// 			errRes := logs.ToAPIErrorCode(err)
// 			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
// 				Nonce: s.NonceManager.GetNonce(),
// 				Error: errRes})
// 			return
// 		}
// 		sr.Address = address

// 		api, err := mefs.NewGateway()
// 		if err != nil {
// 			errRes := logs.ToAPIErrorCode(err)
// 			c.JSON(errRes.HTTPStatusCode, errRes)
// 			return
// 		}

// 		si, err := api.GetPkgSize(c.Request.Context(), address)
// 		if err != nil {
// 			errRes := logs.ToAPIErrorCode(err)
// 			c.JSON(errRes.HTTPStatusCode, errRes)
// 			return
// 		}
// 		sr.StorageList = append(sr.StorageList, si)

// 		api, err = ipfs.NewGateway()
// 		si, err = api.GetPkgSize(c.Request.Context(), address)
// 		if err != nil {
// 			errRes := logs.ToAPIErrorCode(err)
// 			c.JSON(errRes.HTTPStatusCode, errRes)
// 			return
// 		}
// 		sr.StorageList = append(sr.StorageList, si)

// 		c.JSON(http.StatusOK, sr)
// 	})
// }
