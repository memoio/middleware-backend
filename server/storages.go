package server

import (
	"bytes"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/controller"
	"github.com/memoio/backend/internal/logs"
)

func (s Server) StorageRegistRoutes(r *gin.RouterGroup) {
	s.PutobjectRoute(r)
	s.GetObjectRoute(r)
	s.ListObjectsRoute(r)
	s.DeleteObejectRoute(r)
	s.GetObjectPublicRoute(r)
}

func (s Server) PutobjectRoute(r *gin.RouterGroup) {
	s.Router.MaxMultipartMemory = 8 << 20 // 8 MiB

	p := r.Group("/")

	p.POST("/", auth.VerifyIdentityHandler, func(c *gin.Context) {
		address := c.GetString("address")
		chain := c.GetInt("chainid")
		file, err := c.FormFile("file")
		if err != nil {
			errRes := logs.ToAPIErrorCode(logs.ServerError{Message: err.Error()})
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		var public bool
		publics := c.PostForm("public")
		if publics == "true" {
			public = false
		} else {
			public = true
		}
		log.Println(public)

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

		result, err := s.Controller.PutObject(c.Request.Context(), chain, address, object, fr, controller.ObjectOptions{Size: size, UserDefined: ud, Public: public})
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, result)
	})
}

func (s Server) GetObjectRoute(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/:cid", auth.VerifyIdentityHandler, func(c *gin.Context) {
		cid := c.Param("cid")
		address := c.GetString("address")
		chain := c.GetInt("chainid")
		var w bytes.Buffer
		result, err := s.Controller.GetObject(c.Request.Context(), chain, address, cid, &w, controller.ObjectOptions{})
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
	})
}

func (s Server) ListObjectsRoute(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/listobjects", auth.VerifyIdentityHandler, func(c *gin.Context) {
		address := c.GetString("address")
		chain := c.GetInt("chainid")
		result, err := s.Controller.ListObjects(c.Request.Context(), chain, address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, result)
	})
}

func (s Server) DeleteObejectRoute(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/delete", auth.VerifyIdentityHandler, func(c *gin.Context) {
		address := c.GetString("address")
		mid := c.Query("mid")
		chain := c.GetInt("chainid")

		err := s.Controller.DeleteObject(c.Request.Context(), chain, address, mid)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, gin.H{"state": "success"})
	})
}

func (s Server) GetObjectPublicRoute(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/public/:cid", func(c *gin.Context) {
		cid := c.Param("cid")
		chain := c.Query("chainid")
		chainid := big.NewInt(0)
		chainid.SetString(chain, 10)
		var w bytes.Buffer
		result, err := s.Controller.GetObjectPublic(c.Request.Context(), int(chainid.Int64()), cid, &w, controller.ObjectOptions{})
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
	})
}
