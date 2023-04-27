package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/gateway"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
)

func (s Server) mefsRegistRoutes(r *gin.RouterGroup) {
	s.addMefsPutobjectRoutes(r)
	s.addMefsGetObjectRoutes(r)
	s.addMefsListObjectRoutes(r)
	// s.addDeleteRoutes(r)
}

func (s Server) addMefsPutobjectRoutes(r *gin.RouterGroup) {
	s.Router.MaxMultipartMemory = 8 << 20 // 8 MiB
	p := r.Group("/")

	p.POST("/", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			errRes := logs.ToAPIErrorCode(logs.ServerError{err.Error()})
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		if file == nil {
			errRes := logs.ToAPIErrorCode(logs.ServerError{"file is nil"})
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		size := file.Size

		object := file.Filename
		ud := make(map[string]string)

		r, err := file.Open()
		if err != nil {
			errRes := logs.ToAPIErrorCode(logs.ServerError{"open file error"})
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		api, err := mefs.NewGateway()
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		obi, err := gateway.PutObject(c.Request.Context(), api, address, object, r, gateway.ObjectOptions{Size: size, UserDefined: ud})
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		result := make(map[string]string)
		result["id"] = obi.Cid
		c.JSON(http.StatusOK, result)
	})
}

func (s Server) addMefsGetObjectRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/:cid", func(c *gin.Context) {
		cid := c.Param("cid")
		api, err := mefs.NewGateway()
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		obi, err := api.GetObjectInfo(c.Request.Context(), cid)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		var w bytes.Buffer
		err = api.GetObject(c.Request.Context(), cid, &w, gateway.ObjectOptions{})
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		head := fmt.Sprintf("attachment; filename=\"%s\"", obi.Name)
		extraHeaders := map[string]string{
			"Content-Disposition": head,
		}
		c.DataFromReader(http.StatusOK, obi.Size, obi.CType, &w, extraHeaders)
	})
}

func (s Server) addMefsListObjectRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/listobjects", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		api, err := mefs.NewGateway()
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		loi, err := api.ListObjects(c.Request.Context(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		lresponse := ListObjectsResponse{
			Address: address,
			Storage: storage.MEFS.String(),
		}

		for _, oi := range loi {
			lresponse.Object = append(lresponse.Object, ObjectResponse{
				Name:        oi.Name,
				Size:        oi.Size,
				Mid:         oi.Cid,
				ModTime:     oi.ModTime,
				UserDefined: oi.UserDefined,
			})
		}

		c.JSON(http.StatusOK, lresponse)
	})
}
