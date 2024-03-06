package da

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
)

var daStore api.IGateway
var logger = logs.Logger("da")
var defaultDABucket string = "da-bucket"

func LoadDAModule(g *gin.RouterGroup) {
	ui := api.USerInfo{
		Api:   config.Cfg.Storage.Mefs.Api,
		Token: config.Cfg.Storage.Mefs.Token,
	}
	store, err := mefs.NewGatewayWith(ui)
	if err != nil {
		logger.Error("init da-store-mefs error:", err)
		return
	}
	daStore = store

	g.GET("/getObject", getObjectHandler)
	g.POST("/putObject", putObjectHandler)
	g.GET("/warmup", warmupHandler)
}

func getObjectHandler(c *gin.Context) {
	objID := c.Query("id")
	if len(objID) == 0 {
		lerr := logs.ServerError{Message: "object's id is not set"}
		c.Error(lerr)
		return
	}

	var w bytes.Buffer
	err := daStore.GetObject(c.Request.Context(), objID, &w, api.ObjectOptions{})
	if err != nil {
		c.Error(err)
		return
	}

	c.Data(http.StatusOK, utils.TypeByExtension(""), w.Bytes())
}

func putObjectHandler(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)
	data, ok := body["data"].(string)
	if !ok {
		lerr := logs.ServerError{Message: "field 'data' is not set"}
		c.Error(lerr)
		return
	}

	dataHash := crypto.Keccak256Hash([]byte(data))
	var buf bytes.Buffer
	buf.Write([]byte(data))
	oi, err := daStore.PutObject(c.Request.Context(), defaultDABucket, dataHash.Hex(), &buf, api.ObjectOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			c.Error(err)
			return
		}
		// get cid
		cid, err := daStore.(*mefs.Mefs).GetObjectEtag(c.Request.Context(), defaultDABucket, dataHash.Hex())
		if err != nil {
			c.Error(err)
			return
		}
		oi.Cid = cid
	}

	c.JSON(http.StatusOK, gin.H{
		"mid": oi.Cid,
	})
}

func warmupHandler(c *gin.Context) {
	tempStore := daStore.(*mefs.Mefs)
	err := tempStore.MakeBucketWithLocation(c.Request.Context(), defaultDABucket)
	if err != nil {
		if !strings.Contains(err.Error(), "already exist") {
			c.Error(err)
		}
	} else {
		logger.Info("Create bucket ", defaultDABucket)
		for !tempStore.CheckBucket(c.Request.Context(), defaultDABucket) {
			time.Sleep(5 * time.Second)
		}
	}
	c.JSON(http.StatusOK, nil)
}
