package server

import (
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/controller"
	"github.com/memoio/backend/internal/logs"
)

func toInt64(s string) int64 {
	b := new(big.Int)
	b.SetString(s, 10)
	return b.Int64()
}

func (s Server) packagesRegistRoutes(r *gin.RouterGroup) {
	s.addBuyPkgRoutes(r)
	s.addGetPkgListRoutes(r)
	s.addGetBuyPkgRoutes(r)
}

func (s Server) addBuyPkgRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/buypkg", auth.VerifyIdentityHandler, func(c *gin.Context) {
		amount := c.Query("amount")
		pkgid := c.Query("pkgid")
		chainId := c.GetInt("chainid")
		times := time.Now()
		address := c.GetString("address")
		pkg := controller.Package{
			Pkgid:     uint64(toInt64(pkgid)),
			Amount:    toInt64(amount),
			Starttime: uint64(times.Unix()),
			Chainid:   big.NewInt(int64(chainId)).String(),
		}
		receipt, err := s.Controller.BuyPackage(chainId, address, pkg)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		c.JSON(http.StatusOK, receipt)
	})
}

func (s Server) addGetPkgListRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/pkginfos", auth.VerifyIdentityHandler, func(c *gin.Context) {
		chain := c.GetInt("chainid")
		result, err := s.Controller.GetPackageList(chain)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		c.JSON(http.StatusOK, result)
	})
}

func (s Server) addGetBuyPkgRoutes(r *gin.RouterGroup) {
	p := r.Group("/")
	p.GET("/getbuypkgs", auth.VerifyIdentityHandler, func(c *gin.Context) {
		address := c.GetString("address")
		chain := c.GetInt("chainid")
		pi, err := s.Controller.GetUserBuyPackages(chain, address)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		c.JSON(http.StatusOK, pi)
	})
}
