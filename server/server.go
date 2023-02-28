package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/gateway"
)

type Server struct {
	Router  *gin.Engine
	Gateway *gateway.Gateway
	Config  *config.Config
	NonceManager *NonceManager
}

type AuthenticationFaileMessage struct {
	Nonce string
	Error gateway.APIError
}

func NewServer() *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Server")
	})

	nonceManager := NewNonceManager(30 * int64(time.Second.Seconds()), 1 * int64(time.Minute.Seconds()))

	router.GET("/getnonce", func(c *gin.Context) {
		nonce := nonceManager.GetNonce()
		c.String(http.StatusOK, nonce)
	})
	
	router.POST("/login", func(c *gin.Context) {
		var request LoginRequest
		err := c.BindJSON(&request)
		if err != nil{
	    	apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
	        c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(), 
				Error: apiErr, 
			})
			return
		}
		accessToken, freshToken, err := Login(nonceManager, request)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(), 
				Error: apiErr, 
			})
			return
		}

		// if address is new user in "memo.io" {
		// 	init usr info
		// }
		fmt.Println(request.Address)

		c.JSON(http.StatusOK, map[string]string{
			"access token": accessToken, 
			"fresh token": freshToken,
		})
	})

	router.GET("/fresh", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		accessToken, err := VerifyFreshToken(tokenString)
		if err != nil {
			c.String(http.StatusUnauthorized, "Illegal fresh token")
			return
		}
		c.JSON(http.StatusOK, map[string]string{
			"access token": accessToken, 
		})
	})

	config, err := config.ReadFile("")
	if err != nil {
		log.Fatal("config not right")
		return nil
	}
	g := gateway.NewGateway(config)

	s := &Server{
		Router:  router,
		Gateway: g,
		Config:  config,
		NonceManager:  nonceManager, 
	}

	s.registRoute()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.Router,
	}

	return srv
}

func (s Server) registRoute() {
	mefs := s.Router.Group("/mefs")
	s.commonregistRoutes(mefs, gateway.MEFS)
	ipfs := s.Router.Group("/ipfs")
	s.commonregistRoutes(ipfs, gateway.IPFS)
}

func (s Server) commonregistRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	s.addPutobjectRoutes(r, storage)
	s.addGetObjectRoutes(r, storage)
	s.addListObjectRoutes(r, storage)
	s.addGetPriceRoutes(r, storage)
	s.addGetStorageRoutes(r, storage)
	s.addGetBalanceRoutes(r, storage)
	s.addPayRoutes(r, storage)
	s.addS3GetObjectRoutes(r, storage)
}

func (s Server) addPutobjectRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	s.Router.MaxMultipartMemory = 8 << 20 // 8 MiB
	p := r.Group("/")

	p.POST("/", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		file, _ := c.FormFile("file")
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
		obi, err := s.Gateway.PutObject(c.Request.Context(), address, object, r, storage, gateway.ObjectOptions{UserDefined: ud})
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
		balance, err := s.Gateway.GetBalanceInfo(c.Request.Context(), address, storage)
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

func (s Server) addPayRoutes(r *gin.RouterGroup, storage gateway.StorageType) {
	p := r.Group("/")
	p.GET("/pay", func(c *gin.Context) {
		years, ok := c.GetQuery("years")
		if !ok {
			years = "1"
		}
		c.JSON(http.StatusOK, years)
	})
}
