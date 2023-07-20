package routes

import (
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
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

func (h handler) getSpace(c *gin.Context) {
	address := c.GetString("address")
	space, err := h.controller.SpacePayInfo(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, space)
}

func (h handler) getTraffic(c *gin.Context) {
	address := c.GetString("address")
	space, err := h.controller.TrafficPayInfo(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, space)
}
