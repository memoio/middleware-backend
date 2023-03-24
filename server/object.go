package server

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/gateway"
)

func (s Server) addPutobjectRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
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
		result["cid"] = obi.Cid
		c.JSON(http.StatusOK, result)
	})
}

func (s Server) addGetObjectRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	p := r.Group("/")
	p.GET("/:cid", func(c *gin.Context) {
		cid := c.Param("cid")

		if cid == "listobjects" || cid == "balance" || cid == "storage" {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), gateway.AddressError{"address is null"}), gateway.AddressError{"address is null"})
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

func (s Server) addListObjectRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
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

		for _, oi := range loi.Objects {
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

func (s Server) addGetPriceRoutes(r *gin.RouterGroup, stroage gateway.StorageType) {
	p := r.Group("/")
	p.GET("/getprice", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
}

func (s Server) addGetBalanceRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	p := r.Group("/")
	p.GET("/balance", func(c *gin.Context) {
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
		balance, err := s.Gateway.GetBalanceInfo(c.Request.Context(), address)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, apiErr)
			return
		}
		c.JSON(http.StatusOK, BalanceResponse{Address: address, Balance: balance})
	})
}

func (s Server) addGetStorageRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	p := r.Group("/")
	p.GET("/storage", func(c *gin.Context) {
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
		si, err := s.Gateway.GetStorageInfo(c.Request.Context(), address)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, apiErr)
		}
		c.JSON(http.StatusOK, si)
	})
}

func (s Server) addS3GetObjectRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	p := r.Group("/")
	p.GET("/S3/*url", func(c *gin.Context) {
		url := c.Param("url")
		c.JSON(http.StatusOK, url)
	})
}
