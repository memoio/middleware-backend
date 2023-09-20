package auth

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/go-did/memo"
)

func LoadAuthModule(g *gin.RouterGroup) {
	initJWTConfig()

	g.POST("/login", LoginHandler)
	g.GET("/login", GetSessionHandler)
	g.POST("/refresh", RefreshHandler)

	// test API
	g.GET("/test/identity", VerifyIdentityHandler, func(c *gin.Context) {
		c.JSON(200, fmt.Sprintf("did:%s  payload:%s\n", c.GetString("did"), c.GetString("payload")))
	})
	g.GET("/test/accesstoken", VerifyAccessTokenHandler, func(c *gin.Context) {
		c.JSON(200, fmt.Sprintf("did:%s  address:%s\n", c.GetString("did"), c.GetString("address")))
	})
}

// Login godoc
//
//	@Summary		Login
//	@Description	Login
//	@Tags			Login
//	@Accept			json
//	@Produce		json
//	@Param			b	body		string	true	"body"
//	@Success		200	{object}	string	"Login"
//	@Failure		521	{object}	logs.APIError
//	@Failure		400	{object}	logs.APIError
//	@Router			/login [post]
func LoginHandler(c *gin.Context) {
	// var request Request
	body := make(map[string]interface{})
	c.BindJSON(&body)

	did, ok1 := body["did"].(string)
	nonce, ok2 := body["nonce"].(string)
	timestamp, ok3 := body["timestamp"].(float64)
	signature, ok4 := body["signature"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		c.JSON(401, gin.H{"error": "Missing parameters, please refer to the API documentation for details"})
		return
	}

	accessToken, refreshToken, err := Login(did, nonce, int64(timestamp), signature)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})

}

// GetSession godoc
//
//	@Summary		GetSession
//	@Description	GetSession
//	@Tags			GetSession
//	@Accept			json
//	@Produce		json
//	@Param			did	query		string	true	"did"
//	@Success		200	{object}	string	"request id"
//	@Failure		521	{object}	logs.APIError
//	@Failure		400	{object}	logs.APIError
//	@Router			/login{cid} [get]
func GetSessionHandler(c *gin.Context) {
	did := c.Query("did")

	session, err := sessionStore.GetSession(did)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, session)
}

// Refresh godoc
//
//	@Summary		Refresh
//	@Description	Refresh
//	@Tags			Refresh
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization"
//	@Success		200				{object}	string	"Refresh"
//	@Failure		521				{object}	logs.APIError
//	@Failure		400				{object}	logs.APIError
//	@Router			/refresh [post]
func RefreshHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	accessToken, err := VerifyRefreshToken(tokenString)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.AbortWithStatusJSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(200, map[string]string{
		"accessToken": accessToken,
	})

}

func VerifyAccessTokenHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	did, err := VerifyAccessToken(tokenString)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.AbortWithStatusJSON(errRes.HTTPStatusCode, errRes)
		return
	}

	resolver, _ := memo.NewMemoDIDResolver("dev")
	address, _ := resolver.GetMasterKey(did)

	c.Set("address", address)
	c.Set("did", did)
}

func VerifyIdentityHandler(c *gin.Context) {
	ctype := c.GetHeader("Content-Type")
	var did, token, signature, hash string
	var requestID float64
	if ctype == "application/json" {
		body := make(map[string]interface{})
		c.BindJSON(&body)

		did, _ = body["did"].(string)
		token, _ = body["token"].(string)
		requestID, _ = body["requestID"].(float64)
		signature, _ = body["signature"].(string)
		if body["hash"] != nil {
			hash = body["hash"].(string)
		}
	} else {
		did = c.PostForm("did")
		token = c.PostForm("token")
		requestIDStr := c.PostForm("requestID")
		signature = c.PostForm("signature")
		hash = c.PostForm("hash")

		var err error
		requestID, err = strconv.ParseFloat(requestIDStr, 64)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
	}

	if did == "" || token == "" || requestID == 0 || signature == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "Missing parameters, please refer to the API documentation for details"})
		return
	}

	ok, err := VerifyIdentity(did, token, hash, int64(requestID), signature)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.AbortWithStatusJSON(401, gin.H{"error": fmt.Sprintf("failed to verify identity: %s", did)})
		return
	}

	resolver, _ := memo.NewMemoDIDResolver("dev")
	address, _ := resolver.GetMasterKey(did)

	c.Set("address", address)
	c.Set("did", did)
	c.Set("hash", hash)
}
