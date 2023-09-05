package routes

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/server/routes/controller"
)

// storage

// putOBJ godoc
// @Summary		put object
// @Description	put object
// @Tags			OBJ
// @Accept			json
// @Produce		json
// @Param			file	formData	file	true	"Object file to upload"
// @Success		200		{object}	string	"file id"
// @Failure		521		{object}	logs.APIError
// @Router			/mefs/putOBJ/ [post]
// @Router			/ipfs/putOBJ/ [post]
func (h handler) putObjectHandle(c *gin.Context) {
	address := c.GetString("address")
	file, err := c.FormFile("file")
	if err != nil {
		err = logs.ServerError{Message: err.Error()}
		c.Error(err)
	}

	if file == nil {
		err = logs.ServerError{Message: "file is nil"}
		c.Error(err)
		return
	}

	size := file.Size

	object := file.Filename
	ud := make(map[string]string)

	fr, err := file.Open()
	if err != nil {
		err = logs.ServerError{Message: "open file error"}
		c.Error(err)
	}

	sign := c.PostForm("sign")
	area := c.PostForm("area")

	if sign == "" {
		lerr := logs.ServerError{Message: "sign is empty"}
		c.Error(lerr)
	}

	result, err := h.controller.PutObject(c.Request.Context(), address, object, fr, controller.ObjectOptions{Size: size, UserDefined: ud, Sign: sign, Area: area})
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, result)
}

func (h handler) getObjectHandle(c *gin.Context) {
	cid := c.Param("cid")
	address := c.GetString("address")

	sign := c.Query("sign")

	if sign == "" {
		lerr := logs.ServerError{Message: "sign is empty"}
		c.Error(lerr)
	}

	var w bytes.Buffer
	result, err := h.controller.GetObject(c.Request.Context(), address, cid, &w, controller.ObjectOptions{Sign: sign})
	if err != nil {
		c.Error(err)
	}

	head := fmt.Sprintf("attachment; filename=\"%s\"", result.Name)
	extraHeaders := map[string]string{
		"Content-Disposition": head,
	}

	c.DataFromReader(http.StatusOK, result.Size, result.CType, &w, extraHeaders)
}

func (h handler) listObjectsHandle(c *gin.Context) {
	address := c.GetString("address")

	result, err := h.controller.ListObjects(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, result)
}

func (h handler) deleteObjectHandle(c *gin.Context) {
	address := c.GetString("address")
	id := c.Query("id")

	err := h.controller.DeleteObject(c.Request.Context(), address, int(toInt64(id)))
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, gin.H{"state": "success"})
}

// account

func (h handler) getBalanceHandle(c *gin.Context) {
	address := c.GetString("address")
	balance, err := h.controller.GetBalance(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, gin.H{"Address": address, "Balance": balance.String()})
}

func (h handler) getSpaceInfoHandle(c *gin.Context) {
	address := c.GetString("address")
	space, err := h.controller.SpacePayInfo(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, space)
}

func (h handler) getTrafficInfoHandle(c *gin.Context) {
	address := c.GetString("address")
	pi, err := h.controller.TrafficPayInfo(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, pi)
}

func (h handler) getSpaceCheckHashHandle(c *gin.Context) {
	address := c.GetString("address")
	filesize := c.Query("size")

	res, err := h.controller.GetSpaceCheckHash(c.Request.Context(), address, toUint64(filesize))
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) getTrafficCheckHashHandle(c *gin.Context) {
	address := c.GetString("address")
	filesize := c.Query("size")

	res, err := h.controller.GetTrafficCheckHash(c.Request.Context(), address, toUint64(filesize))
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, res)
}

// spacePrice godoc
//
//	@Summary		get space price
//	@Description	get space price
//	@Tags			price
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	int	"file id"
//	@Failure		521	{object}	logs.APIError
//	@Router			/account/getSpacePrice/ [get]
func (h handler) getSpacePriceHandle(c *gin.Context) {
	price, err := h.controller.GetSpacePrice(c.Request.Context())
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, price)
}

func (h handler) getTrafficPriceHandle(c *gin.Context) {
	res, err := h.controller.GetTrafficPrice(c.Request.Context())
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) buySpaceHandle(c *gin.Context) {
	address := c.GetString("address")
	size := c.Query("size")

	res, err := h.controller.BuySpace(c.Request.Context(), address, toUint64(size))
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) buyTrafficHandle(c *gin.Context) {
	address := c.GetString("address")
	checksize := c.Query("size")

	res, err := h.controller.BuyTraffic(c.Request.Context(), address, toUint64(checksize))
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) getApproveTsHash(c *gin.Context) {
	address := c.GetString("address")
	value := c.Query("value")
	pt := c.Query("type")

	res, err := h.controller.Approve(c.Request.Context(), api.StringToPayType(pt), address, toBigInt(value))
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) getAllowanceHandle(c *gin.Context) {
	address := c.GetString("address")
	at := c.Query("type")

	res, err := h.controller.Allowance(c.Request.Context(), api.StringToPayType(at), address)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) checkReceiptHandle(c *gin.Context) {
	receipt := c.Query("receipt")

	err := h.controller.CheckReceipt(c.Request.Context(), receipt)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, "success")
}

func (h handler) cashSpaceHandle(c *gin.Context) {
	address := c.Query("address")

	res, err := h.controller.CashSpace(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

func (h handler) cashTrafficHandle(c *gin.Context) {
	address := c.Query("address")

	res, err := h.controller.CashTraffic(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}
