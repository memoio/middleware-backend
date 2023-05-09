package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/controller"
	"github.com/memoio/backend/internal/logs"
)

// func (s Server) getAddress(c *gin.Context) (string, error) {
// 	tokenString := c.GetHeader("Authorization")
// 	address, err := VerifyAccessToken(tokenString)
// 	if err != nil {
// 		errRes := logs.ToAPIErrorCode(err)
// 		c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
// 			Nonce: s.NonceManager.GetNonce(),
// 			Error: errRes})
// 		return "", err
// 	}

// 	return address, nil
// }

func (s Server) StorageRegistRoutes(r *gin.RouterGroup) {
	s.PutobjectRoute(r)
	s.GetObjectRoute(r)
	s.ListObjectsRoute(r)
}

func (s Server) PutobjectRoute(r *gin.RouterGroup) {
	s.Router.MaxMultipartMemory = 8 << 20 // 8 MiB

	p := r.Group("/")

	p.POST("/", func(c *gin.Context) {
		// address, err := s.getAddress(c)
		// if err != nil {
		// 	return
		// }
		address := c.Query("address")

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

		result, err := controller.PutObject(c.Request.Context(), r.BasePath(), address, object, fr, controller.ObjectOptions{Size: size, UserDefined: ud})
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
	p.GET("/:cid", func(c *gin.Context) {
		cid := c.Param("cid")

		var w bytes.Buffer
		result, err := controller.GetObject(c.Request.Context(), r.BasePath(), cid, &w, controller.ObjectOptions{})
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
	p.GET("/listobjects", func(c *gin.Context) {
		// address, err := s.getAddress(c)
		// if err != nil {
		// 	return
		// }
		address := c.Query("address")

		result, err := controller.ListObjects(c.Request.Context(), r.BasePath(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, result)
	})
}
