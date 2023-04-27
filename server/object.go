package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/storage"
)

func (s Server) addGetPriceRoutes(r *gin.RouterGroup, stroage storage.StorageType) {
	p := r.Group("/")
	p.GET("/getprice", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
}

// func (s Server) addDeleteRoutes(r *gin.RouterGroup) {
// 	p := r.Group("/")
// 	p.GET("/delete", func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")
// 		address, err := VerifyAccessToken(tokenString)
// 		if err != nil {
// 			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
// 			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
// 				Nonce: s.NonceManager.GetNonce(),
// 				Error: apiErr,
// 			})
// 			return
// 		}
// 		mid := c.Query("mid")
// 		err = s.Gateway.MefsDeleteObject(c.Request.Context(), address, mid)
// 		if err != nil {
// 			c.JSON(521, err.Error())
// 			return
// 		}
// 		c.JSON(http.StatusOK, response{Status: "Success"})
// 	})
// }

func (s Server) addS3GetObjectRoutes(r *gin.RouterGroup, storage storage.StorageType) {
	p := r.Group("/")
	p.GET("/S3/*url", func(c *gin.Context) {
		url := c.Param("url")
		c.JSON(http.StatusOK, url)
	})
}
