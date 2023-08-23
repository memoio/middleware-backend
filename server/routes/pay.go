package routes

import (
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

func toInt64(s string) int64 {
	b := new(big.Int)
	b.SetString(s, 10)
	return b.Int64()
}

func toUint64(s string) uint64 {
	b := new(big.Int)
	b.SetString(s, 10)
	return b.Uint64()
}

func toBigInt(s string) *big.Int {
	b := new(big.Int)
	b.SetString(s, 10)
	return b
}
func (h handler) getBalanceHandle(c *gin.Context) {
	address := c.GetString("did")
	balance, err := h.controller.GetBalance(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Address": address, "Balance": balance.String()})
}

func (h handler) getSpace(c *gin.Context) {
	address := c.GetString("did")
	space, err := h.controller.SpacePayInfo(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, space)
}

func (h handler) getTraffic(c *gin.Context) {
	address := c.GetString("did")
	space, err := h.controller.TrafficPayInfo(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, space)
}

func (h handler) cashTraffic(c *gin.Context) {
	address := c.Query("address")
	res, err := h.controller.CashTraffic(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}

// spacePrice godoc
//	@Summary		get space price
//	@Description	get space price
//	@Tags			price
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	int	"file id"
//	@Failure		521	{object}	logs.APIError
//	@Router			/mefs/getSpacePrice/ [get]
//	@Router			/ipfs/getSpacePrice/ [get]
func (h handler) spacePrice(c *gin.Context) {
	res, err := h.controller.SpacePrice(c.Request.Context())
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) trafficPrice(c *gin.Context) {
	res, err := h.controller.TrafficPrice(c.Request.Context())
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) BuySpace(c *gin.Context) {
	address := c.GetString("did")
	checksize := c.Query("size")

	res, err := h.controller.BuySpace(c.Request.Context(), address, toUint64(checksize))
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) BuyTraffic(c *gin.Context) {
	address := c.GetString("did")
	checksize := c.Query("size")

	res, err := h.controller.BuyTraffic(c.Request.Context(), address, toUint64(checksize))
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) Approve(c *gin.Context) {
	address := c.GetString("did")
	value := c.Query("value")
	pt := c.Query("type")

	res, err := h.controller.Approve(c.Request.Context(), api.StringToPayType(pt), address, toBigInt(value))
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) allowance(c *gin.Context) {
	address := c.GetString("did")
	at := c.Query("type")

	res, err := h.controller.Allowance(c.Request.Context(), api.StringToPayType(at), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}
func (h handler) cashSpace(c *gin.Context) {
	address := c.Query("address")
	res, err := h.controller.CashSpace(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) spaceHash(c *gin.Context) {
	address := c.GetString("did")
	filesize := c.Query("size")

	res, err := h.controller.GetStorePayHash(c.Request.Context(), address, toUint64(filesize))
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) trafficHash(c *gin.Context) {
	address := c.GetString("did")
	filesize := c.Query("size")

	res, err := h.controller.GetReadPayHash(c.Request.Context(), address, toUint64(filesize))
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}
	c.JSON(http.StatusOK, res)
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
