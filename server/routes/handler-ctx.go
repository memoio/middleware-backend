package routes

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/server/routes/controller"
	"github.com/memoio/middleware-response/response"
)

func (h *handler) getStore(c *gin.Context) error {
	store, ok := c.Get("store")
	if !ok || store == nil {
		lerr := logs.ServerError{Message: "store not set"}
		c.Error(lerr)
		return lerr
	}

	storei := store.(api.IGateway)
	h.controller.SetStore(storei)

	return nil
}

// storage

// putOBJ godoc
//
//	@Summary		put object
//	@Description	put object
//	@Tags			PutObj
//	@Accept			json
//	@Produce		json
//	@Param			did			formData	string	true	"did"
//	@Param			token		formData	string	true	"token"
//	@Param			requestID	formData	uint64	true	"requestID"
//	@Param			signature	formData	string	true	"signature"
//	@Param			file		formData	file	true	"file"
//	@Param			sign		formData	string	true	"sign"
//	@Param			area		formData	string	false	"area"
//	@Success		200			{object}	string	"file id"
//	@Failure		521			{object}	logs.APIError
//	@Failure		400			{object}	logs.APIError
//	@Failure		525			{object}	logs.APIError
//	@Router			/mefs/putObject/ [post]
//	@Router			/ipfs/putObject/ [post]
func (h handler) putObjectHandle(c *gin.Context) {
	address := c.GetString("address")
	err := h.getStore(c)
	if err != nil {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		err = logs.ServerError{Message: err.Error()}
		c.Error(err)
		return
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
		return
	}

	sign := c.PostForm("sign")
	area := c.PostForm("area")

	if sign == "" {
		lerr := logs.ServerError{Message: "sign is empty"}
		c.Error(lerr)
		return
	}

	result, err := h.controller.PutObject(c.Request.Context(), address, object, fr, controller.ObjectOptions{Size: size, UserDefined: ud, Sign: sign, Area: area})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// getObject godoc
//
//	@Summary		getObject
//	@Description	getObject
//	@Tags			getObject
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"body"
//	@Param			sign	query		string	true	"sign"
//	@Param			cid		path		string	true	"cid"
//	@Success		200		{object}	string	"file id"
//	@Failure		521		{object}	logs.APIError
//	@Failure		400		{object}	logs.APIError
//	@Router			/mefs/getObject/{cid} [post]
//	@Router			/ipfs/getObject/{cid} [post]
func (h handler) getObjectHandle(c *gin.Context) {
	err := h.getStore(c)
	if err != nil {
		return
	}

	cid := c.Param("cid")
	address := c.GetString("address")

	sign := c.Query("sign")

	if sign == "" {
		lerr := logs.ServerError{Message: "sign is empty"}
		c.Error(lerr)
		return
	}

	var w bytes.Buffer
	result, err := h.controller.GetObject(c.Request.Context(), address, cid, &w, controller.ObjectOptions{Sign: sign})
	if err != nil {
		c.Error(err)
		return
	}

	head := fmt.Sprintf("attachment; filename=\"%s\"", result.Name)
	extraHeaders := map[string]string{
		"Content-Disposition": head,
	}

	c.DataFromReader(http.StatusOK, result.Size, result.CType, &w, extraHeaders)
}

// listObjects godoc
//
//	@Summary		listObjects
//	@Description	listObjects
//	@Tags			listObjects
//	@Accept			json
//	@Produce		json
//	@Param			b	body		string	true	"body"
//	@Success		200	{object}	string	"objs"
//	@Failure		521	{object}	logs.APIError
//	@Failure		400	{object}	logs.APIError
//	@Router			/mefs/listObject/ [post]
//	@Router			/ipfs/listObject/ [post]
func (h handler) listObjectsHandle(c *gin.Context) {
	address := c.GetString("address")

	result, err := h.controller.ListObjects(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// deleteObjec godoc
//
//	@Summary		deleteObjec
//	@Description	deleteObjec
//	@Tags			deleteObjec
//	@Accept			json
//	@Produce		json
//	@Param			b	body		string	true	"body"
//	@Param			id	query		string	true	"id"
//	@Success		200	{object}	string	"file id"
//	@Failure		521	{object}	logs.APIError
//	@Failure		400	{object}	logs.APIError
//	@Router			/mefs/deleteObject [post]
//	@Router			/ipfs/deleteObject [post]
func (h handler) deleteObjectHandle(c *gin.Context) {
	err := h.getStore(c)
	if err != nil {
		return
	}

	address := c.GetString("address")
	id := c.Query("id")

	err = h.controller.DeleteObject(c.Request.Context(), address, int(toInt64(id)))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"state": "success"})
}

// account

// getBalance godoc
//
//	@Summary		getBalance
//	@Description	getBalance
//	@Tags			getBalance
//	@Accept			json
//	@Produce		json
//	@Param			b	body		string	true	"b"
//	@Success		200	{object}	int		"balance"
//	@Failure		521	{object}	logs.APIError
//	@Router			/mefs/getBalance [post]
func (h handler) getBalanceHandle(c *gin.Context) {
	address := c.GetString("address")
	balance, err := h.controller.GetBalance(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
		return
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
//	@Router			/mefs/getSpaceInfo [post]
func (h handler) getSpaceInfoHandle(c *gin.Context) {
	address := c.GetString("address")
	space, err := h.controller.SpacePayInfo(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, space)
}

// getTrafficInfo godoc
//
//	@Summary		getTrafficInfo
//	@Description	getTrafficInfo
//	@Tags			getTrafficInfo
//	@Accept			json
//	@Produce		json
//	@Param			b	body		string	true	"b"
//	@Success		200	{object}	int		"getTrafficInfo"
//	@Failure		521	{object}	logs.APIError
//	@Router			/mefs/getTrafficInfo [post]
func (h handler) getTrafficInfoHandle(c *gin.Context) {
	address := c.GetString("address")
	pi, err := h.controller.TrafficPayInfo(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pi)
}

// getSpaceCheck godoc
//
//	@Summary		getSpaceCheck
//	@Description	getSpaceCheck
//	@Tags			getSpaceCheck
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			size	query		string	true	"size"
//	@Success		200		{object}	int		"getSpaceCheck"
//	@Failure		521		{object}	logs.APIError
//	@Router			/mefs/getSpaceCheck [post]
func (h handler) getSpaceCheckHandle(c *gin.Context) {
	address := c.Query("address")
	filesize := c.Query("size")

	res, err := h.controller.GetSpaceCheckHash(c.Request.Context(), address, toUint64(filesize))
	if err != nil {
		c.Error(err)
		return
	}
	response := response.CheckResponse(res)
	logger.Info(response.Hash())
	c.JSON(http.StatusOK, response)
}

// getTrafficCheck godoc
//
//	@Summary		getTrafficCheck
//	@Description	getTrafficCheck
//	@Tags			getTrafficCheck
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			size	query		string	true	"size"
//	@Success		200		{object}	int		"getTrafficCheck"
//	@Failure		521		{object}	logs.APIError
//	@Router			/mefs/getTrafficCheck [post]
func (h handler) getTrafficCheckHandle(c *gin.Context) {
	address := c.GetString("address")
	filesize := c.Query("size")

	res, err := h.controller.GetTrafficCheckHash(c.Request.Context(), address, toUint64(filesize))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// getSpacePrice godoc
//
//	@Summary		getSpacePrice
//	@Description	getSpacePrice
//	@Tags			getSpacePrice
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	int	"file id"
//	@Failure		521	{object}	logs.APIError
//	@Router			/mefs/getSpacePrice/ [get]
func (h handler) getSpacePriceHandle(c *gin.Context) {
	price, err := h.controller.GetSpacePrice(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, price)
}

// getTrafficPrice godoc
//
//	@Summary		getTrafficPrice
//	@Description	getTrafficPrice
//	@Tags			getTrafficPrice
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	int	"getTrafficPrice"
//	@Failure		521	{object}	logs.APIError
//	@Router			/mefs/getTrafficPrice/ [get]
func (h handler) getTrafficPriceHandle(c *gin.Context) {
	res, err := h.controller.GetTrafficPrice(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// spaceHash godoc
//
//	@Summary		BuySpace
//	@Description	get buy space tx hash
//	@Tags			BuySpace
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			size	query		string	true	"size"
//	@Success		200		{object}	int		"BuySpace"
//	@Failure		521		{object}	logs.APIError
//	@Router			/mefs/buySpace [post]
func (h handler) buySpaceHandle(c *gin.Context) {
	address := c.GetString("address")
	size := c.Query("size")

	res, err := h.controller.BuySpace(c.Request.Context(), address, toUint64(size))
	if err != nil {
		c.Error(err)
		return
	}
	rsp := response.Transaction(res)

	reps, err := rsp.Marshal()
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		c.Error(lerr)
		return
	}
	c.JSON(http.StatusOK, hexutil.Encode(reps))
}

// BuyTraffic godoc
//
//	@Summary		BuyTraffic
//	@Description	get buy traffic tx hash
//	@Tags			BuyTraffic
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			size	query		string	true	"size"
//	@Success		200		{object}	int		"BuyTraffic"
//	@Failure		521		{object}	logs.APIError
//	@Router			/mefs/buyTraffic [post]
func (h handler) buyTrafficHandle(c *gin.Context) {
	address := c.GetString("address")
	checksize := c.Query("size")

	res, err := h.controller.BuyTraffic(c.Request.Context(), address, toUint64(checksize))
	if err != nil {
		c.Error(err)
		return
	}
	rsp := response.Transaction(res)

	reps, err := rsp.Marshal()
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		c.Error(lerr)
		return
	}
	c.JSON(http.StatusOK, hexutil.Encode(reps))
}

// Approve godoc
//
//	@Summary		recharge
//	@Description	recharge
//	@Tags			recharge
//	@Accept			json
//	@Produce		json
//	@Param			b		body		string	true	"b"
//	@Param			value	query		string	true	"value"
//	@Param			type	query		string	true	"type"
//	@Success		200		{object}	int		"getApproveTsHash"
//	@Failure		521		{object}	logs.APIError
//	@Router			/mefs/recharge [post]
func (h handler) getApproveTsHash(c *gin.Context) {
	address := c.GetString("address")
	value := c.Query("value")
	pt := c.Query("type")

	res, err := h.controller.Approve(c.Request.Context(), api.StringToPayType(pt), address, toBigInt(value))
	if err != nil {
		c.Error(err)
		return
	}
	rsp := response.Transaction(res)

	reps, err := rsp.Marshal()
	if err != nil {
		lerr := logs.ControllerError{Message: err.Error()}
		c.Error(lerr)
		return
	}
	c.JSON(http.StatusOK, hexutil.Encode(reps))
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
//	@Router			/mefs/getAllowance [post]
func (h handler) getAllowanceHandle(c *gin.Context) {
	address := c.GetString("address")
	at := c.Query("type")

	res, err := h.controller.Allowance(c.Request.Context(), api.StringToPayType(at), address)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// checkReceipt godoc
//
//	@Summary		checkReceipt
//	@Description	checkReceipt
//	@Tags			checkReceipt
//	@Accept			json
//	@Produce		json
//	@Param			receipt	query		string	true	"receipt"
//	@Success		200		{object}	int		"cashSpace"
//	@Failure		521		{object}	logs.APIError
//	@Router			/mefs/getReceipt [get]
func (h handler) checkReceiptHandle(c *gin.Context) {
	receipt := c.Query("receipt")

	err := h.controller.CheckReceipt(c.Request.Context(), receipt)
	if err != nil {
		c.Error(err)
		return
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
//	@Router			/mefs/cashSpace [get]
func (h handler) cashSpaceHandle(c *gin.Context) {
	address := c.Query("address")

	res, err := h.controller.CashSpace(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// cashTraffic godoc
//
//	@Summary		cashTraffic
//	@Description	cashTraffic
//	@Tags			cashTraffic
//	@Accept			json
//	@Produce		json
//	@Param			address	query		string	true	"address"
//	@Success		200		{object}	int		"cashTraffic"
//	@Failure		521		{object}	logs.APIError
//	@Router			/mefs/cashTraffic [get]
func (h handler) cashTrafficHandle(c *gin.Context) {
	address := c.Query("address")

	res, err := h.controller.CashTraffic(c.Request.Context(), address)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, res)
}
