package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/gateway"
)

func LoginHandler(nonceManager *NonceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request LoginRequest
		err := c.BindJSON(&request)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		accessToken, freshToken, err := Login(nonceManager, request)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}

		// if address is new user in "memo.io" {
		// 	init usr info
		// }
		fmt.Println(request.Address)

		c.JSON(http.StatusOK, map[string]string{
			"access token": accessToken,
			"fresh token":  freshToken,
		})
	}
}

func LensLoginHandler(nonceManager *NonceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request EIP4361Request
		err := c.BindJSON(&request)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		accessToken, freshToken, err := LoginWithMethod(nonceManager, request, LensMod)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}

		// if address is new user in "memo.io" {
		// 	init usr info
		// }
		// fmt.Println(request.Address)

		c.JSON(http.StatusOK, map[string]string{
			"access token": accessToken,
			"fresh token":  freshToken,
		})
	}
}

func FreshHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		accessToken, err := VerifyFreshToken(tokenString)
		if err != nil {
			c.String(http.StatusUnauthorized, "Illegal fresh token")
			return
		}
		c.JSON(http.StatusOK, map[string]string{
			"access token": accessToken,
		})
	}
}
