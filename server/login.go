package server

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/logs"
)

func ChallengeHandler(nonceManager *NonceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		uri, err := url.Parse(c.GetHeader("Origin"))
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: errRes})
			return
		}
		domain := uri.Host
		nonce := nonceManager.GetNonce()

		challenge, err := Challenge(domain, address, uri.String(), nonce)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: errRes})
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
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: errRes})
			return
		}
		accessToken, refreshToken, _, err := Login(nonceManager, request)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: errRes})
			return
		}

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
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: errRes})
			return
		}
		accessToken, refreshToken, _, isRegistered, err := LoginWithLens(request, checkRegistered)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: errRes})
			return
		}

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

		c.JSON(http.StatusOK, map[string]string{
			"accessToken": accessToken,
		})
	}
}
