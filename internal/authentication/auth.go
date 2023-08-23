package auth

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/logs"
)

func LoadAuthModule(g *gin.RouterGroup, checkRegistered bool) {
	InitRecommendTable()

	g.GET("/challenge", ChallengeHandler())

	g.POST("/login", LoginHandler())

	g.POST("/lens/login", LensLoginHandler(checkRegistered))

	g.GET("/refresh", RefreshHandler())

	g.GET("/identity", VerifyIdentityHandler, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"address": c.GetString("address"),
			"chainid": c.GetInt("chainid"),
		})
	})

	g.GET("/recommend", ListRecommendHandler())

	g.GET("/recommend/:address", GetRecommendHandler())
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
		accessToken, refreshToken, address, err := Login(nonceManager, request)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		var newAccount = true
		var recommend = Recommend{
			Address:     address,
			Recommender: request.Recommender,
			Source:      request.Source,
		}
		err = recommend.CreateRecommend()
		if err != nil {
			if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
				errRes := logs.ToAPIErrorCode(err)
				c.JSON(errRes.HTTPStatusCode, errRes)
				return
			} else {
				newAccount = false
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"newAccount":   newAccount,
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
	if tokenString == "" {
		tokenString = "Bearer " + c.Query("token")
	}

	address, chainid, userID, err := VerifyAccessToken(tokenString)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.AbortWithStatusJSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.Set("address", address)
	c.Set("chainid", chainid)
	c.Set("userid", userID)
}

func ListRecommendHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		recommender := c.Query("recommender")
		recommends, err := ListRecommend(recommender)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.AbortWithStatusJSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, recommends)
	}
}

func GetRecommendHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Param("address")
		recommend, err := GetRecommend(address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.AbortWithStatusJSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, recommend)
	}
}
