package server

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/gateway"
	db "github.com/memoio/backend/global/database"
)

func ChallengeHandler(nonceManager *NonceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		uri, err := url.Parse(c.GetHeader("Origin"))
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		domain := uri.Host
		nonce := nonceManager.GetNonce()

		challenge, err := Challenge(domain, address, uri.String(), nonce)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		c.String(http.StatusOK, challenge)
	}
}

func LoginHandler(nonceManager *NonceManager) gin.HandlerFunc {
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
		accessToken, refreshToken, err := Login(nonceManager, request)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		address, err := VerifyAccessToken(accessToken)
		if err != nil {
			c.String(http.StatusUnauthorized, "Illegal Lens token")
			return
		}
		db.AddressInfo{Address: address}.Insert()
		// if address is new user in "memo.io" {
		// 	init usr info
		// }
		// fmt.Println(request.Address)

		c.JSON(http.StatusOK, map[string]string{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

func LensLoginHandler(nonceManager *NonceManager, checkRegistered bool) gin.HandlerFunc {
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
		accessToken, refreshToken, isRegistered, err := LoginWithLens(request, checkRegistered)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		address, err := VerifyAccessToken(accessToken)
		if err != nil {
			c.String(http.StatusUnauthorized, "Illegal Lens token")
			return
		}
		db.AddressInfo{Address: address}.Insert()
		// if address is new user in "memo.io" {
		// 	init usr info
		// }
		// fmt.Println(request.Address)

		c.JSON(http.StatusOK, map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"isRegistered": isRegistered,
		})

	}
}

func RefreshHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		accessToken, err := VerifyRefreshToken(tokenString)
		if err != nil {
			c.String(http.StatusUnauthorized, "Illegal refresh token")
			return
		}

		// address, err := VerifyAccessToken(tokenString)
		// if err != nil {
		// 	c.String(http.StatusUnauthorized, "Illegal refresh token")
		// 	return
		// }

		// db.AddressInfo{Address: address}.Insert()
		c.JSON(http.StatusOK, map[string]string{
			"accessToken": accessToken,
		})
	}
}
