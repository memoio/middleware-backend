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
//
//	@Summary		put object
//	@Description	put object
//	@Tags			OBJ
//	@Accept			json
//	@Produce		json
//	@Param			file	formData	file	true	"Object file to upload"
//	@Success		200		{object}	string	"file id"
//	@Failure		521		{object}	logs.APIError
//	@Router			/mefs/putOBJ/ [post]
//	@Router			/ipfs/putOBJ/ [post]
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

// get balance godoc
//
//	@Summary		get balance
//	@Description	get balance
//	@Tags			balance
//	@Accept			json
//	@Produce		json
//	@Param			b	body		string	true	"b"
//	@Success		200	{object}	int		"balance"
//	@Failure		521	{object}	logs.APIError
//	@Router			/account/getBalance [post]
func (h handler) getBalanceHandle(c *gin.Context) {
	address := c.GetString("address")
	balance, err := h.controller.GetBalance(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, gin.H{"Address": address, "Balance": balance.String()})
}

// getSpaceInfo godoc
//
//	@Summary		getSpaceInfo
//	@Description	getSpaceInfo
//	@Tags			getSpaceInfo
//	@Accept			json
//	@Produce		json
//	@Param			b	body		string	true	"b"
//	@Success		200	{object}	int		"getSpaceInfo"
//	@Failure		521	{object}	logs.APIError
//	@Router			/account/getSpaceInfo [post]
func (h handler) getSpaceInfoHandle(c *gin.Context) {
	address := c.GetString("address")
	space, err := h.controller.SpacePayInfo(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, space)
}

// getTrafficInfoHandle godoc
//
//	@Summary		getTrafficInfoHandle
//	@Description	getTrafficInfoHandle
//	@Tags			getTrafficInfoHandle
//	@Accept			json
//	@Produce		json
//	@Param			b	body		string	true	"b"
//	@Success		200	{object}	int		"getTrafficInfoHandle"
//	@Failure		521	{object}	logs.APIError
//	@Router			/account/getTrafficInfo [post]
func (h handler) getTrafficInfoHandle(c *gin.Context) {
	address := c.GetString("address")
	pi, err := h.controller.TrafficPayInfo(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, pi)
}

// spaceHash godoc
//
//	@Summary		spaceHash
//	@Description	spaceHash
//	@Tags			spaceHash
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			size	query		string	true	"size"
//	@Success		200		{object}	int		"spaceHash"
//	@Failure		521		{object}	logs.APIError
//	@Router			/account/getSpaceHash [post]
func (h handler) getSpaceCheckHashHandle(c *gin.Context) {
	address := c.GetString("address")
	filesize := c.Query("size")

	res, err := h.controller.GetSpaceCheckHash(c.Request.Context(), address, toUint64(filesize))
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

// trafficHash godoc
//
//	@Summary		trafficHash
//	@Description	trafficHash
//	@Tags			trafficHash
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			size	query		string	true	"size"
//	@Success		200		{object}	int		"trafficHash"
//	@Failure		521		{object}	logs.APIError
//	@Router			/account/getTrafficHash [post]
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

// trafficPrice godoc
//
//	@Summary		get trafficPrice
//	@Description	get trafficPrice
//	@Tags			traffic price
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	int	"getTrafficPrice"
//	@Failure		521	{object}	logs.APIError
//	@Router			/account/getTrafficPrice/ [get]
func (h handler) getTrafficPriceHandle(c *gin.Context) {
	res, err := h.controller.GetTrafficPrice(c.Request.Context())
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

// spaceHash godoc
//
//	@Summary		BuySpace
//	@Description	BuySpace
//	@Tags			BuySpace
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			size	query		string	true	"size"
//	@Success		200		{object}	int		"BuySpace"
//	@Failure		521		{object}	logs.APIError
//	@Router			/account/buySpace [post]
func (h handler) buySpaceHandle(c *gin.Context) {
	address := c.GetString("address")
	size := c.Query("size")

	res, err := h.controller.BuySpace(c.Request.Context(), address, toUint64(size))
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

// BuyTraffic godoc
//
//	@Summary		BuyTraffic
//	@Description	BuyTraffic
//	@Tags			BuyTraffic
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			size	query		string	true	"size"
//	@Success		200		{object}	int		"BuyTraffic"
//	@Failure		521		{object}	logs.APIError
//	@Router			/account/buyTraffic [post]
func (h handler) buyTrafficHandle(c *gin.Context) {
	address := c.GetString("address")
	checksize := c.Query("size")

	res, err := h.controller.BuyTraffic(c.Request.Context(), address, toUint64(checksize))
	if err != nil {
		c.Error(err)
	}
	c.JSON(http.StatusOK, res)
}

// Approve godoc
//
//	@Summary		getApproveTsHash
//	@Description	getApproveTsHash
//	@Tags			getApproveTsHash
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			value	query		string	true	"value"
//	@Param			type	query		string	true	"type"
//	@Success		200		{object}	int		"getApproveTsHash"
//	@Failure		521		{object}	logs.APIError
//	@Router			/account/getApproveTsHash [post]
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

// allowance godoc
//
//	@Summary		getAllowance
//	@Description	getAllowance
//	@Tags			getAllowance
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			type	query		string	true	"type"
//	@Success		200		{object}	int		"getAllowance"
//	@Failure		521		{object}	logs.APIError
//	@Router			/account/getAllowance [post]
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

// cashSpace godoc
//
//	@Summary		cashSpace
//	@Description	cashSpace
//	@Tags			cashSpace
//	@Accept			json
//	@Produce		json
//	@Param			address	query		string	true	"address"
//	@Success		200		{object}	int		"cashSpace"
//	@Failure		521		{object}	logs.APIError
//	@Router			/admin/cashSpace [get]
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
