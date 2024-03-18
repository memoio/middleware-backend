package da

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/kzg"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
	proof "github.com/memoio/go-did/file-proof"
)

var DefaultSRS *kzg.SRS
var zeroCommit bls12381.G1Affine
var zeroProof kzg.OpeningProof
var defaultProofInstance *proof.ProofInstance
var userSk *ecdsa.PrivateKey
var submitterSk *ecdsa.PrivateKey
var daStore api.IGateway
var logger = logs.Logger("da")
var defaultDABucket string = "da-bucket"
var defaultDAObject string = "da-txdata"
var defaultExpiration time.Duration = 7 * 24 * time.Hour

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

	DefaultSRS, err = kzg.NewSRS(4*1024, big.NewInt(985))
	if err != nil {
		logger.Error("init da-SRS error:", err)
		return
	}

	// poly := split(nil)
	// commit, err := kzg.Commit(poly, DefaultSRS.Pk)
	// if err != nil {
	// 	logger.Error("load zero-commit error:", err)
	// 	return
	// }
	// zeroCommit = commit
	// logger.Info(zeroCommit)

	zeroCommit.X.SetZero()
	zeroCommit.Y.SetZero()

	zeroProof.ClaimedValue.SetZero()
	zeroProof.H.X.SetZero()
	zeroProof.H.Y.SetZero()

	err = InitDAFileInfoTable()
	if err != nil {
		logger.Error("init da-file-info error:", err)
		return
	}

	userSk, err = crypto.HexToECDSA(config.Cfg.DataAccess.UserSecurityKey)
	if err != nil {
		logger.Error("generate user private key error:", err)
		return
	}

	submitterSk, err = crypto.HexToECDSA(config.Cfg.DataAccess.SubmitterSecurityKey)
	if err != nil {
		logger.Error("generate submitter private key error:", err)
		return
	}

	defaultProofInstance, err = proof.NewProofInstance(userSk, "dev")
	if err != nil {
		logger.Error("init proof instance error:", err)
		return
	}

	g.GET("/getObject", getObjectHandler)
	g.POST("/putObject", putObjectHandler)
	g.GET("/warmup", warmupHandler)
	fmt.Println("load da moudle success!")
}

func getObjectHandler(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		lerr := logs.ServerError{Message: "object's id is not set"}
		c.Error(lerr)
		return
	}

	var commit bls12381.G1Affine
	idBytes, err := hexutil.Decode(id)
	if err != nil {
		c.Error(err)
		return
	}
	err = commit.Unmarshal(idBytes)
	if err != nil {
		c.Error(err)
		return
	}

	file, err := GetFileInfoByCommit(commit)
	if err != nil {
		c.Error(err)
		return
	}

	var w bytes.Buffer
	err = daStore.GetObject(c.Request.Context(), file.Mid, &w, api.ObjectOptions{})
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

	databyte, err := hex.DecodeString(data)
	if err != nil {
		lerr := logs.ServerError{Message: "field 'data' is not legally hexadecimal presented"}
		c.Error(lerr)
		return
	}

	object := defaultDAObject + hex.EncodeToString(crypto.Keccak256(databyte))

	var buf *bytes.Buffer = bytes.NewBuffer(databyte)
	oi, err := daStore.PutObject(c.Request.Context(), defaultDABucket, object, buf, api.ObjectOptions{})
	if err != nil {
		c.Error(err)
		return
	}

	elements := split(databyte)
	// log.Println(string(buf.Bytes()), elements)
	commit, err := kzg.Commit(elements, DefaultSRS.Pk)
	if err != nil {
		c.Error(err)
		return
	}

	start := time.Now()
	end := start.Add(defaultExpiration)
	hash := defaultProofInstance.GetCredentialHash(commit, uint64(oi.Size), big.NewInt(start.Unix()), big.NewInt(end.Unix()))
	signature, err := crypto.Sign(hash, submitterSk)
	if err != nil {
		c.Error(err)
		return
	}

	err = defaultProofInstance.AddFile(commit, uint64(oi.Size), big.NewInt(start.Unix()), big.NewInt(end.Unix()), signature)
	if err != nil {
		c.Error(err)
		return
	}

	var fileInfo = DAFileInfo{
		Commit:     commit,
		Mid:        oi.Cid,
		Size:       oi.Size,
		Expiration: end.Unix(),
	}

	err = fileInfo.CreateDAFileInfo()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": hexutil.Encode(commit.Marshal()),
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
