package routes

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/server/controller"
)

func (h handler) putObjectHandle(c *gin.Context) {
	address := c.GetString("did")
	file, err := c.FormFile("file")
	if err != nil {
		errRes := logs.ToAPIErrorCode(logs.ServerError{Message: err.Error()})
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
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

	sign := c.PostForm("sign")
	area := c.PostForm("area")

	if sign == "" {
		lerr := logs.ControllerError{Message: "sign is empty"}
		errRes := logs.ToAPIErrorCode(lerr)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	result, err := h.controller.PutObject(c.Request.Context(), address, object, fr, controller.ObjectOptions{Size: size, UserDefined: ud, Sign: sign, Area: area})
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h handler) getObjectHandle(c *gin.Context) {
	cid := c.Param("cid")
	address := c.GetString("did")

	sign := c.Query("sign")

	if sign == "" {
		lerr := logs.ControllerError{Message: "sign is empty"}
		errRes := logs.ToAPIErrorCode(lerr)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	var w bytes.Buffer
	result, err := h.controller.GetObject(c.Request.Context(), address, cid, &w, controller.ObjectOptions{Sign: sign})
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

func (h handler) listObjectsHandle(c *gin.Context) {
	address := c.GetString("did")

	result, err := h.controller.ListObjects(c.Request.Context(), address)
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h handler) deleteObjectHandle(c *gin.Context) {
	address := c.GetString("did")
	id := c.Query("id")

	err := h.controller.DeleteObject(c.Request.Context(), address, int(toInt64(id)))
	if err != nil {
		errRes := logs.ToAPIErrorCode(err)
		c.JSON(errRes.HTTPStatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, gin.H{"state": "success"})
}
