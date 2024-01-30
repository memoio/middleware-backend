package routes

import (
	"bytes"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/config"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/controller"
	"github.com/memoio/backend/internal/gateway/ipfs"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/utils"
)

type handler struct {
	controller *controller.Controller
}

func handlerMefs() handler {
	store, err := mefs.NewGateway()
	if err != nil {
		logger.Error("init mefs error:", err)
	}

	config, err := config.ReadFile()
	if err != nil {
		logger.Error("config not right ", err)

	}
	control, err := controller.NewController(api.MEFS, store, config)
	if err != nil {
		logger.Error("get control error:", err)
	}
	return handler{controller: control}
}

func handlerIpfs() handler {
	store, err := ipfs.NewGateway()
	if err != nil {
		logger.Error("init ipfs error:", err)
	}
	config, err := config.ReadFile()
	if err != nil {
		logger.Error("config not right ", err)

	}

	control, err := controller.NewController(api.IPFS, store, config)
	if err != nil {
		logger.Error("get control error:", err)
	}
	return handler{controller: control}
}

func handleStorage(r *gin.RouterGroup, h handler) {
	r.POST("/", auth.VerifyIdentityHandler, h.putObjectHandle)
	r.GET("/:cid", auth.VerifyIdentityHandler, h.getObjectHandle)
	r.GET("/listobjects", auth.VerifyIdentityHandler, h.listObjectsHandle)
	r.GET("/delete", auth.VerifyIdentityHandler, h.deleteObjectHandle)
	r.GET("/public/:cid", h.getPublicObjectHandle)

	r.GET("/pkginfos", auth.VerifyIdentityHandler, h.getPackageInfoHandle)
	r.GET("/buypkg", auth.VerifyIdentityHandler, h.buyPackageHandle)
	r.GET("/getbuypkgs", auth.VerifyIdentityHandler, h.getBuyPackagesHandle)
	r.GET("/receipt", h.receiptHandle)

	r.GET("/balance", auth.VerifyIdentityHandler, h.getBalanceHandle)
	r.GET("/storageinfo", auth.VerifyIdentityHandler, h.getStorageInfoHandle)
	r.GET("/flowsize", auth.VerifyIdentityHandler, h.getFlowSizeHandle)
}

func (h handler) putObjectHandle(c *gin.Context) {
	address := c.GetString("address")
	chain := c.GetInt("chainid")
	user := c.PostForm("user")
	file, err := c.FormFile("file")
	if err != nil {
		errRes := logs.ToAPIErrorCode(logs.ServerError{Message: err.Error()})
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	key := c.PostForm("key")
	if key == "" {
		key = "f1d4a0b37124c3a7"
	}

	var public bool
	publics := c.PostForm("public")
	if publics == "true" {
		key = ""
		public = true
	} else {
		public = false
	}

	if file == nil {
		errRes := logs.ToAPIErrorCode(logs.ServerError{Message: "file is nil"})
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	size := file.Size

	object := file.Filename
	ud := make(map[string]string)

	fr, err := file.Open()
	if err != nil {
		errRes := logs.ToAPIErrorCode(logs.ServerError{Message: "open file error"})
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	if key != "" {
		re, err := utils.EncryptFile(fr, []byte(key))
		if err != nil {
			lerr := logs.ControllerError{Message: fmt.Sprint("encryt error", err)}
			errRes := logs.ToAPIErrorCode(lerr)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		result, err := h.controller.PutObject(c.Request.Context(), chain, address, object, re, controller.ObjectOptions{Size: size, UserDefined: ud, Public: public})
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		c.JSON(http.StatusOK, result)
		return
	}
	result, err := h.controller.PutObject(c.Request.Context(), chain, address, object, fr, controller.ObjectOptions{Size: size, UserDefined: ud, Public: public, User: user})
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h handler) getObjectHandle(c *gin.Context) {
	cid := c.Param("cid")
	address := c.GetString("address")
	chain := c.GetInt("chainid")
	var w bytes.Buffer

	key := c.PostForm("key")

	result, err := h.controller.GetObject(c.Request.Context(), chain, address, cid, &w, controller.ObjectOptions{Key: []byte(key)})
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	if !result.Public {
		if key == "" {
			key = "f1d4a0b37124c3a7"
		}
		output := new(bytes.Buffer)
		output.Write(w.Bytes())
		w.Reset()
		err = utils.DecryptFile(output, &w, []byte(key))
		if err != nil {
			lerr := logs.ControllerError{Message: fmt.Sprint("encryt error", err)}
			errRes := logs.ToAPIErrorCode(lerr)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
	}

	head := fmt.Sprintf("attachment; filename=\"%s\"", result.Name)
	extraHeaders := map[string]string{
		"Content-Disposition": head,
	}

	c.DataFromReader(http.StatusOK, result.Size, result.CType, &w, extraHeaders)
}

func (h handler) listObjectsHandle(c *gin.Context) {
	address := c.GetString("address")
	chain := c.GetInt("chainid")
	result, err := h.controller.ListObjects(c.Request.Context(), chain, address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h handler) deleteObjectHandle(c *gin.Context) {
	address := c.GetString("address")
	id := c.Query("id")
	if id == "" {
		msg := logs.ControllerError{Message: "id not set, please check storage id"}
		errRes := logs.ToAPIErrorCode(msg)
		c.JSON(errRes.HTTPStatusCode, errRes )
		return
	}
	err := h.controller.DeleteObject(c.Request.Context(), address, int(toInt64(id)))
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, gin.H{"state": "success"})
}

func (h handler) getPublicObjectHandle(c *gin.Context) {
	cid := c.Param("cid")
	chain := c.Query("chainid")
	chainid := big.NewInt(0)
	chainid.SetString(chain, 10)
	var w bytes.Buffer
	result, err := h.controller.GetObjectPublic(c.Request.Context(), int(chainid.Int64()), cid, &w, controller.ObjectOptions{})
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	head := fmt.Sprintf("attachment; filename=\"%s\"", result.Name)
	extraHeaders := map[string]string{
		"Content-Disposition": head,
	}

	c.DataFromReader(http.StatusOK, result.Size, result.CType, &w, extraHeaders)
}

func (h handler) getPackageInfoHandle(c *gin.Context) {
	chain := c.GetInt("chainid")
	result, err := h.controller.GetPackageList(chain)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h handler) buyPackageHandle(c *gin.Context) {
	amount := c.Query("amount")
	pkgid := c.Query("pkgid")
	chainId := c.GetInt("chainid")
	times := time.Now()
	address := c.GetString("address")
	pkg := controller.Package{
		Pkgid:     uint64(toInt64(pkgid)),
		Amount:    toInt64(amount),
		Starttime: uint64(times.Unix()),
		Chainid:   big.NewInt(int64(chainId)).String(),
	}
	receipt, err := h.controller.BuyPackage(chainId, address, pkg)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, receipt)
}

func (h handler) getBuyPackagesHandle(c *gin.Context) {
	address := c.GetString("address")
	chain := c.GetInt("chainid")
	pi, err := h.controller.GetUserBuyPackages(chain, address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, pi)
}

func (h handler) receiptHandle(c *gin.Context) {
	receipt := c.Query("receipt")
	chain := c.GetInt("chainid")
	err := h.controller.CheckReceipt(c.Request.Context(), chain, receipt)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h handler) getBalanceHandle(c *gin.Context) {
	address := c.GetString("address")
	chain := c.GetInt("chainid")
	balance, err := h.controller.GetBalance(c.Request.Context(), chain, address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Address": address, "Balance": balance.String()})
}

func (h handler) getStorageInfoHandle(c *gin.Context) {
	address := c.GetString("address")
	chain := c.GetInt("chainid")

	si, err := h.controller.GetStorageInfo(c.Request.Context(), chain, address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, si)
}

func (h handler) getFlowSizeHandle(c *gin.Context) {
	address := c.GetString("address")
	chain := c.GetInt("chainid")

	res, err := h.controller.GetFlowSize(c.Request.Context(), chain, address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, res)
}
