package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LoadAuthRouter(r *gin.RouterGroup) {
	r.POST("/login", LoginHandler)
	r.GET("/identity", VerifyIdentityHandler, func(c *gin.Context) {
		c.JSON(200, fmt.Sprintf("address:%s  chainid:%d\n", c.GetString("address"), c.GetInt64("chainid")))
	})
}

func LoginHandler(c *gin.Context) {
	// var request Request
	body := make(map[string]interface{})
	c.BindJSON(&body)

	address, ok1 := body["address"].(string)
	token, ok2 := body["token"].(string)
	timestamp, ok3 := body["timestamp"].(float64)
	signature, ok4 := body["signature"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		c.JSON(401, gin.H{"error": "Missing parameters, please refer to the API documentation for details"})
		return
	}

	chainID, ok := body["chainid"].(float64)
	if !ok {
		chainID = 985
	}

	_, err := Login(address, token, int64(chainID), int64(timestamp), signature)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.String(200, "Login success!")
}

func VerifyIdentityHandler(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)

	token, ok1 := body["token"].(string)
	requestID, ok2 := body["requestID"].(float64)
	signature, ok3 := body["signature"].(string)
	if !ok1 || !ok2 || !ok3 {
		c.AbortWithStatusJSON(401, gin.H{"error": "Missing parameters, please refer to the API documentation for details"})
		return
	}

	chainID, ok := body["chainid"].(float64)
	if !ok {
		chainID = 985
	}

	address, err := VerifyIdentity(token, int64(chainID), int64(requestID), signature)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
		return
	}

	c.Set("address", address)
	c.Set("chainid", int64(chainID))
}
