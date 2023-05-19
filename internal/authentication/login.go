package auth

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/logs"
)

func LoadAuthModule(g *gin.RouterGroup, checkRegistered bool) {
	g.GET("/challenge", ChallengeHandler())

	g.POST("/login", LoginHandler())

	g.POST("/lens/login", LensLoginHandler(checkRegistered))

	g.GET("/refresh", RefreshHandler())
}

func ChallengeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		uri, err := url.Parse(c.GetHeader("Origin"))
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		domain := uri.Host
		nonce := nonceManager.GetNonce()

		var chainID int
		if c.Query("chainid") != "" {
			chainID, err = strconv.Atoi(c.Query("chainid"))
			if err != nil {
				errRes := logs.ToAPIErrorCode(err)
				c.JSON(errRes.HTTPStatusCode, errRes)
				return
			}
		} else {
			chainID = 985
		}

		challenge, err := Challenge(domain, address, uri.String(), nonce, chainID)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		c.String(http.StatusOK, challenge)
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request EIP4361Request
		err := c.BindJSON(&request)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		accessToken, refreshToken, _, err := Login(nonceManager, request)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
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

func LensLoginHandler(checkRegistered bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request EIP4361Request
		err := c.BindJSON(&request)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		accessToken, refreshToken, _, isRegistered, err := LoginWithLens(request, checkRegistered)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
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

func VerifyIdentityHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	address, err := VerifyAccessToken(tokenString)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.AbortWithStatusJSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.Set("address", address)
}
