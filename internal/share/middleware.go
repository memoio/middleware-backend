package share

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/api"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/gateway/ipfs"
	"github.com/memoio/backend/internal/gateway/mefs"
	"github.com/memoio/backend/internal/logs"
	"github.com/memoio/backend/internal/storage"
	"github.com/memoio/backend/utils"
)

var ApiMap map[string]*Api

type Api struct {
	G api.IGateway
	T storage.StorageType
}

func init() {
	loadApiMap()
}

func loadApiMap() {
	ApiMap = make(map[string]*Api)

	mefs, err := mefs.NewGateway()
	if err != nil {
		log.Println("load mefs ap failed")
		return
	}
	ApiMap["/mefs"] = &Api{G: mefs, T: storage.MEFS}

	ipfs, err := ipfs.NewGateway()
	if err != nil {
		log.Println("load ipfs ap failed")
		return
	}
	ApiMap["/ipfs"] = &Api{G: ipfs, T: storage.IPFS}
}

func LoadShareModule(g *gin.RouterGroup) {
	err := InitShareTable()
	if err != nil {
		panic(err.Error())
	}

	{
		// 免费
		share := g.Group("share", ShareAvailableHandler())

		// 下载分享文件
		share.GET("/:shareid", BeforeDownloadHandler(), DownloadShareHandler())

		// 获取分享信息
		share.GET("info/:shareid", GetShareHandler())
	}

	{
		// 需要登录
		share := g.Group("share", auth.VerifyIdentityHandler)

		// 创建分享
		share.POST("", CreateShareHandler())

		// 列出分享
		share.GET("", ListSharesHandler())

		// 将分享添加到我的文件列表中
		share.POST("save/:shareid", ShareAvailableHandler(), BeforeDownloadHandler(), SaveShareHandler())

		// 删除分享
		share.DELETE(":shareid", ShareAvailableHandler(), DeleteShareHandler())
	}
}

func ShareAvailableHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		share := GetShareByID(c.Param("shareid"))

		if share == nil || !share.IsAvailable() {
			c.AbortWithStatusJSON(404, "The share link is not available")
			return
		}

		c.Set("share", share)
		// c.Next()
	}
}

func BeforeDownloadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// password := c.Query("password")
		address := c.GetString("address")
		chainID := c.GetInt("chainid")

		shareObj, _ := c.Get("share")
		share := shareObj.(*ShareObjectInfo)

		err := share.CanDownload(address, chainID)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.AbortWithStatusJSON(errRes.HTTPStatusCode, errRes)
			return
		}

		// c.Next()
	}
}

func DownloadShareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		shareObj, _ := c.Get("share")
		share := shareObj.(*ShareObjectInfo)

		file, err := share.Source()
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		var w bytes.Buffer
		err = ApiMap["/"+share.SType.String()].G.GetObject(c.Request.Context(), file.Mid, &w, api.ObjectOptions{})
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		head := fmt.Sprintf("attachment; filename=\"%s\"", file.Name)
		extraHeaders := map[string]string{
			"Content-Disposition": head,
		}

		c.DataFromReader(http.StatusOK, file.Size, utils.TypeByExtension(file.Name), &w, extraHeaders)

	}
}

func CreateShareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.GetString("address")
		chainID := c.GetInt("chainid")

		var request CreateShareRequest
		err := c.ShouldBindJSON(&request)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		res, err := CreateShare(address, chainID, request)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}
		c.JSON(200, res)
	}
}

func GetShareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// password := c.Query("password")
		address := c.GetString("address")
		chainID := c.GetInt("chainid")

		shareObj, _ := c.Get("share")
		share := shareObj.(*ShareObjectInfo)

		share, err := GetShare(address, chainID, share)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, share)
	}
}

func SaveShareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.GetString("address")
		chainID := c.GetInt("chainid")

		shareObj, _ := c.Get("share")
		share := shareObj.(*ShareObjectInfo)

		err := SaveShare(address, chainID, share)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, "add share success")
	}
}

func DeleteShareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.GetString("address")
		chainID := c.GetInt("chainid")

		shareObj, _ := c.Get("share")
		share := shareObj.(*ShareObjectInfo)

		err := DeleteShare(address, chainID, share)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, "delete success")
	}
}

func ListSharesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.GetString("address")
		chainID := c.GetInt("chainid")

		shares, err := ListShares(address, chainID)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		c.JSON(http.StatusOK, shares)
	}
}
