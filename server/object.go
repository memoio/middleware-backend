package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/gateway"
	"github.com/memoio/backend/internal/storage"
)

func (s Server) addPutobjectRoutes(r *gin.RouterGroup, storage storage.StorageType) {
	s.Router.MaxMultipartMemory = 8 << 20 // 8 MiB
	p := r.Group("/")

	p.POST("/", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		file, _ := c.FormFile("file")
		size := file.Size

		object := file.Filename
		ud := make(map[string]string)
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: s.NonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		r, err := file.Open()
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, apiErr)
			return
		}
		obi, err := s.Gateway.PutObject(c.Request.Context(), address, object, r, storage, gateway.ObjectOptions{Size: size, UserDefined: ud})
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, apiErr)
			return
		}
		result := make(map[string]string)
		result["id"] = obi.Cid
		c.JSON(http.StatusOK, result)
	})
}

func (s Server) addGetObjectRoutes(r *gin.RouterGroup, storage storage.StorageType) {
	p := r.Group("/")
	p.GET("/:cid", func(c *gin.Context) {
		cid := c.Param("cid")

		if cid == "listobjects" || cid == "balance" || cid == "storage" {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), gateway.AddressError{Message: "address is null"}), gateway.AddressError{Message: "address is null"})
			c.JSON(apiErr.HTTPStatusCode, apiErr)
			return
		}
		obi, err := s.Gateway.GetObjectInfo(c.Request.Context(), storage, cid)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, apiErr)
			return
		}

		var w bytes.Buffer
		err = s.Gateway.GetObject(c.Request.Context(), cid, storage, &w, gateway.ObjectOptions{})
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, apiErr)
			return
		}
		head := fmt.Sprintf("attachment; filename=\"%s\"", obi.Name)
		extraHeaders := map[string]string{
			"Content-Disposition": head,
		}
		c.DataFromReader(200, obi.Size, obi.CType, &w, extraHeaders)
	})
}

func (s Server) addListObjectRoutes(r *gin.RouterGroup, storage storage.StorageType) {
	p := r.Group("/")
	p.GET("/listobjects", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: s.NonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}

		loi, err := s.Gateway.ListObjects(c.Request.Context(), address, storage)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, apiErr)
			return
		}

		lresponse := ListObjectsResponse{
			Address: address,
			Storage: storage.String(),
		}

		for _, oi := range loi {
			lresponse.Object = append(lresponse.Object, ObjectResponse{
				Name:        oi.Name,
				Size:        oi.Size,
				Cid:         oi.Cid,
				ModTime:     oi.ModTime,
				UserDefined: oi.UserDefined,
			})
		}

		c.JSON(http.StatusOK, lresponse)
	})
}

func (s Server) addGetPriceRoutes(r *gin.RouterGroup, stroage storage.StorageType) {
	p := r.Group("/")
	p.GET("/getprice", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
}

func (s Server) addDeleteRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/delete", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: s.NonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		mid := c.Query("mid")
		err = s.Gateway.MefsDeleteObject(c.Request.Context(), address, mid)
		if err != nil {
			c.JSON(521, err.Error())
			return
		}
		c.JSON(http.StatusOK, response{Status: "Success"})
	})
}

func (s Server) addS3GetObjectRoutes(r *gin.RouterGroup, storage storage.StorageType) {
	p := r.Group("/")
	p.GET("/S3/*url", func(c *gin.Context) {
		url := c.Param("url")
		c.JSON(http.StatusOK, url)
	})
}
