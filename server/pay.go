package server

import (
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/controller"
	"github.com/memoio/backend/internal/logs"
)

func toInt64(s string) int64 {
	b := new(big.Int)
	b.SetString(s, 10)
	return b.Int64()
}

type response struct {
	Status string
}

func toResponse(f bool) response {
	if f {
		return response{Status: "Success!"}
	} else {
		return response{Status: "Failed!"}
	}

}

func (s Server) addBuyPkgRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/buypkg", func(c *gin.Context) {
		amount := c.Query("amount")
		pkgid := c.Query("pkgid")
		chainId := c.Query("chainid")
		times := time.Now()
		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: s.NonceManager.GetNonce(),
				Error: errRes})
			return
		}
		pkg := controller.Package{
			Pkgid:     uint64(toInt64(pkgid)),
			Amount:    toInt64(amount),
			Starttime: uint64(times.Second()),
			Chainid:   chainId,
		}
		flag := s.Controller.BuyPackage(address, pkg)
		if !flag {
			c.JSON(521, "buy pkg failed")
		}
		c.JSON(http.StatusOK, toResponse(flag))
	})
}

func (s Server) addGetPkgListRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/pkginfos", func(c *gin.Context) {
		result, err := s.Controller.GetPackageList()
		if err != nil {
			c.JSON(522, err.Error())
		}
		c.JSON(http.StatusOK, result)
	})
}

func (s Server) addGetBuyPkgRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/getbuypkgs", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		address, err := VerifyAccessToken(tokenString)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, AuthenticationFaileMessage{
				Nonce: s.NonceManager.GetNonce(),
				Error: errRes})
			return
		}

		pi, err := s.Controller.GetUserBuyPackages(address)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, pi)
	})
}
