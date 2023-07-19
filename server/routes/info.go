package routes

import (
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

func toInt64(s string) int64 {
	b := new(big.Int)
	b.SetString(s, 10)
	return b.Int64()
}

func (h handler) getBalanceHandle(c *gin.Context) {
	address := c.GetString("address")
	balance, err := h.controller.GetBalance(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Address": address, "Balance": balance.String()})
}

func (h handler) getStorageInfoHandle(c *gin.Context) {
	address := c.GetString("address")

	si, err := h.controller.GetStorageInfo(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, si)
}

func (h handler) getPkgInfos(c *gin.Context) {
	result, err := h.controller.GetPackageList(c.Request.Context())
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h handler) getFlowSize(c *gin.Context) {
	address := c.GetString("address")

	res, err := h.controller.GetFlowSize(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h handler) getBuyPackages(c *gin.Context) {
	address := c.GetString("address")

	pi, err := h.controller.GetUserBuyPackages(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, pi)
}

func (h handler) checkReceipt(c *gin.Context) {
	receipt := c.Query("receipt")

	err := h.controller.CheckReceipt(c.Request.Context(), receipt)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h handler) buyPackage(c *gin.Context) {
	amount := c.Query("amount")
	pkgid := c.Query("pkgid")
	chainId := c.GetInt("chainid")
	times := time.Now()
	address := c.GetString("address")
	pkg := api.BuyPackage{
		Pkgid:     uint64(toInt64(pkgid)),
		Amount:    toInt64(amount),
		Starttime: uint64(times.Unix()),
		Chainid:   big.NewInt(int64(chainId)).String(),
	}
	receipt, err := h.controller.BuyPackage(c.Request.Context(), address, pkg)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, receipt)
}
