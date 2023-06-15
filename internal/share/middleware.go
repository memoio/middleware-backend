package share

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/memoio/backend/internal/authentication"
	"github.com/memoio/backend/internal/logs"
)

func LoadAuthModule(g *gin.RouterGroup) {
	{
		// 免费
		share := g.Group("share", ShareAvailableHandler())

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

		log.Println(c.Param("shareid"))
		log.Println(share)

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

		err = share.DownloadBy(address, chainID)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.AbortWithStatusJSON(errRes.HTTPStatusCode, errRes)
			return
		}

		// c.Next()
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
		address := c.GetString("address")
		chainID := c.GetInt("chainID")

		shareObj, _ := c.Get("share")
		share := shareObj.(*ShareObjectInfo)

		var request GetShareRequest
		err := c.BindJSON(&request)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes)
			return
		}

		share, err = GetShare(address, chainID, share, request)
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
		chainID := c.GetInt("chainID")

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
		shareID := c.Param("shareid")

		err := DeleteShare(address, chainID, shareID)
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
