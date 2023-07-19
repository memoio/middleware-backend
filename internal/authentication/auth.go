package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LoadAuthRouter(r *gin.RouterGroup) {
	r.POST("/login", LoginHandler)
	r.GET("/identity", VerifyIdentityHandler, func(c *gin.Context) {
		c.JSON(200, fmt.Sprintf("did:%s  payload:%s\n", c.GetString("did"), c.GetString("payload")))
	})
}

func LoginHandler(c *gin.Context) {
	// var request Request
	body := make(map[string]interface{})
	c.BindJSON(&body)

	did, ok1 := body["did"].(string)
	token, ok2 := body["token"].(string)
	timestamp, ok3 := body["timestamp"].(float64)
	signature, ok4 := body["signature"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		c.JSON(401, gin.H{"error": "Missing parameters, please refer to the API documentation for details"})
		return
	}

	ok, err := Login(did, token, int64(timestamp), signature)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(401, gin.H{"error": fmt.Sprintf("can't log in as %s", did)})
		return
	}

	c.String(200, "Login success!")
}

func VerifyIdentityHandler(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)

	did, ok1 := body["did"].(string)
	token, ok2 := body["token"].(string)
	requestID, ok3 := body["requestID"].(float64)
	signature, ok4 := body["signature"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		c.AbortWithStatusJSON(401, gin.H{"error": "Missing parameters, please refer to the API documentation for details"})
		return
	}

	var payload string
	if body["payload"] != nil {
		payload = body["payload"].(string)
	}
	ok, err := VerifyIdentity(did, token, payload, int64(requestID), signature)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(401, gin.H{"error": fmt.Sprintf("failed to verify identity: %s", did)})
		return
	}

	// c.Set("address", address)
	c.Set("did", did)
	c.Set("payload", payload)
}
