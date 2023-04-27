package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/gateway/ipfs"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
)

func (s Server) testRegistRoutes(r *gin.RouterGroup) {
	// p := r.Group("/")
	// p.GET("/storage", func(c *gin.Context) {
	// 	address := c.Query("address")
	// 	si, err := s.Gateway.GetPkgSize(c.Request.Context(), storage.MEFS, address)
	// 	if err != nil {
	// 		c.JSON(516, err.Error())
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, si)
	// })

	// p.POST("/put", func(c *gin.Context) {
	// 	address := c.Query("address")

	// 	file, err := c.FormFile("file")
	// 	if err != nil {
	// 		c.JSON(511, fmt.Sprintf("read file err %s", err))
	// 		return
	// 	}
	// 	if file == nil {
	// 		c.JSON(511, "get file error")
	// 		return
	// 	}
	// 	size := file.Size

	// 	object := file.Filename
	// 	ud := make(map[string]string)
	// 	if err != nil {
	// 		apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
	// 		c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
	// 			Nonce: s.NonceManager.GetNonce(),
	// 			Error: apiErr,
	// 		})
	// 		return
	// 	}
	// 	r, err := file.Open()
	// 	if err != nil {
	// 		apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
	// 		c.JSON(apiErr.HTTPStatusCode, apiErr)
	// 		return
	// 	}
	// 	obi, err := s.Gateway.MefsPutObject(c.Request.Context(), address, object, r, gateway.ObjectOptions{Size: size, UserDefined: ud})
	// 	if err != nil {
	// 		apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
	// 		c.JSON(apiErr.HTTPStatusCode, apiErr)
	// 		return
	// 	}
	// 	result := make(map[string]string)
	// 	result["id"] = obi.Cid
	// 	c.JSON(http.StatusOK, result)
	// })

	// p.POST("/update", func(c *gin.Context) {
	// 	address := c.Query("address")
	// 	hashid := c.Query("hash")

	// 	si, err := s.Gateway.TestUpdatePkg(c.Request.Context(), storage.MEFS, address, hashid, 1024)
	// 	if err != nil {
	// 		log.Println("TEST: ", err)
	// 		c.JSON(520, err.Error())
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, si)
	// })
	// p.GET("/delete", func(c *gin.Context) {
	// 	address := "0x2Dc689e597fA3545F0c5f6aF2f4c1De2d334C8EC"
	// 	hashid := c.Query("mid")

	// 	r := contract.StoreOrderPkgExpiration(address, hashid, uint8(storage.MEFS), big.NewInt(1124))
	// 	c.JSON(http.StatusOK, toResponse(r))
	// })

	// p.GET("/balance", func(c *gin.Context) {
	// 	address := c.Query("address")

	// 	balance := contract.BalanceOf(c.Request.Context(), address)
	// 	c.JSON(http.StatusOK, BalanceResponse{Address: address, Balance: balance.String()})
	// })

	// p.GET("/pay", func(c *gin.Context) {
	// 	address := c.Query("address")
	// 	hashid := c.Query("hash")

	// 	p := s.Gateway.TestPay(c.Request.Context(), address, hashid, 1, 1024)

	// 	c.JSON(http.StatusOK, p)
	// })

	// p.GET("/pkginfo", func(c *gin.Context) {
	// 	result, err := contract.StoreGetPkgInfos()
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, err.Error())
	// 	}
	// 	c.JSON(http.StatusOK, result)
	// })
	// p.GET("/setpkg", func(c *gin.Context) {
	// 	time := c.Query("time")
	// 	amount := c.Query("amount")
	// 	kind := c.Query("kind")
	// 	size := c.Query("size")

	// 	flag := contract.AdminAddPkgInfo(time, amount, kind, size)
	// 	c.JSON(http.StatusOK, flag)
	// })
	// p.GET("buypkg", func(c *gin.Context) {
	// 	address := c.Query("address")
	// 	chainId := c.Query("chainid")
	// 	times := time.Now()
	// 	flag := contract.StoreBuyPkg(address, 1, 1, uint64(times.Second()), chainId)
	// 	c.JSON(http.StatusOK, toResponse(flag))
	// })
	// p.GET("getbuypkg", func(c *gin.Context) {
	// 	address := c.Query("address")
	// 	pi, err := contract.StoreGetBuyPkgs(address)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, err.Error())
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, pi)
	// })
	// p.GET("/getall", func(c *gin.Context) {
	// 	a := contract.GetStoreAllSize()
	// 	c.JSON(http.StatusOK, a)
	// })
	// p.GET("/error", func(c *gin.Context) {
	// 	c.JSON(533, logs.ErrResponse{})
	// })
	s.addTestGetStorageRoutes(r)
}

func (s Server) addTestGetStorageRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/getstorage", func(c *gin.Context) {
		var sr StorageResponse

		address := c.Query("address")
		sr.Address = address

		api, err := mefs.NewGateway()
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		si, err := api.GetPkgSize(c.Request.Context(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		sr.StorageList = append(sr.StorageList, si)

		api = ipfs.NewGateway(s.Config.Storage.Ipfs.Host)
		si, err = api.GetPkgSize(c.Request.Context(), address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		sr.StorageList = append(sr.StorageList, si)

		c.JSON(http.StatusOK, sr)
	})
}
