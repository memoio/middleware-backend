package market

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/logs"
)

func LoadNFTMarketModule(g *gin.RouterGroup) {
	g.GET("show", ShowNFTHandler())

	g.GET("search", SearchNFTHandler())
}

func ShowNFTHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.Query("page")
		sizeStr := c.Query("size")
		order := c.Query("order")
		ascend := c.Query("asc")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes.Description)
			return
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes.Description)
			return
		}

		var nftInfos []NFT
		if ascend == "false" {
			nftInfos, err = ListNFT(page, size, order, false)
		} else {
			nftInfos, err = ListNFT(page, size, order, true)
		}

		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes.Description)
			return
		}

		c.JSON(200, nftInfos)
	}
}

func SearchNFTHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.Query("page")
		sizeStr := c.Query("size")
		text := c.Query("query")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes.Description)
			return
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			errRes := logs.ToAPIErrorCode(err)
			c.JSON(errRes.HTTPStatusCode, errRes.Description)
			return
		}

		results := Search(text, page, size)

		var tokenIDs []string = make([]string, len(results))
		for index, result := range results {
			tokenIDs[index] = result.DocId
		}

		var nfts []NFT
		err = database.DataBase.Where("token_id IN ?", tokenIDs).Find(&nfts).Error
		if err != nil {
			c.JSON(524, err.Error())
			return
		}

		for index := range nfts {
			nfts[index].TokenID -= 1
		}

		c.JSON(200, nfts)
	}
}
