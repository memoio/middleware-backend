package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/config"
	"github.com/memoio/backend/gateway"
)

type Server struct {
	Router       *gin.Engine
	Gateway      *gateway.Gateway
	Config       *config.Config
	NonceManager *NonceManager
}

type AuthenticationFaileMessage struct {
	Nonce string
	Error gateway.APIError
}

func NewServer(endpoint string, checkRegistered bool) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(Cors())
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Server")
	})

	nonceManager := NewNonceManager(30*int64(time.Second.Seconds()), 1*int64(time.Minute.Seconds()))

	router.GET("/challenge", func(c *gin.Context) {
		address := c.Query("address")
		uri, err := url.Parse(c.GetHeader("Origin"))
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		domain := uri.Host
		nonce := nonceManager.GetNonce()

		fmt.Println(address, domain, uri, nonce)

		challenge, err := Challenge(domain, address, uri.String(), nonce)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		c.String(http.StatusOK, challenge)
	})

	router.POST("/login", func(c *gin.Context) {
		var request EIP4361Request
		err := c.BindJSON(&request)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		accessToken, refreshToken, err := Login(nonceManager, request)
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
		// fmt.Println(request.Address)

		c.JSON(http.StatusOK, map[string]string{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	})

	router.POST("/lens/login", func(c *gin.Context) {
		var request EIP4361Request
		err := c.BindJSON(&request)
		if err != nil {
			apiErr := gateway.ErrorCodes.ToAPIErrWithErr(gateway.ToAPIErrorCode(c.Request.Context(), err), err)
			c.JSON(apiErr.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: nonceManager.GetNonce(),
				Error: apiErr,
			})
			return
		}
		accessToken, refreshToken, isRegistered, err := LoginWithLens(request, checkRegistered)
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
		// fmt.Println(request.Address)

		c.JSON(http.StatusOK, map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"isRegistered": isRegistered,
		})
	})

	router.GET("/refresh", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		accessToken, err := VerifyRefreshToken(tokenString)
		if err != nil {
			c.String(http.StatusUnauthorized, "Illegal refresh token")
			return
		}
		c.JSON(http.StatusOK, map[string]string{
			"accessToken": accessToken,
		})
	})

	config, err := config.ReadFile("")
	if err != nil {
		log.Fatal("config not right")
		return nil
	}
	InitAuthConfig(config.SecurityKey, config.Domain, config.LensAPIUrl)
	g := gateway.NewGateway(config)

	s := &Server{
		Router:       router,
		Gateway:      g,
		Config:       config,
		NonceManager: nonceManager,
	}

	s.registRoute()

	srv := &http.Server{
		Addr:    endpoint,
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
