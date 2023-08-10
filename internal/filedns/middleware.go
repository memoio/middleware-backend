package filedns

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/memoio/backend/internal/logs"
)

func LoadFileDnsModule(g *gin.RouterGroup) {
	g.GET("/search", SearchHandler())
}

type SearchRespond struct {
	Mid      string
	FileName string
	Keywords []string
	Size     int64
	ModTime  time.Time
}

func SearchHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		text := c.Query("query")
		pageStr := c.Query("page")
		sizeStr := c.Query("size")

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

		ids := Search(text, page-1, size)

		var res []SearchRespond
		for _, id := range ids {
			// get did document info

			// get file info
			// database.DataBase.Where("did = ?", id.DocId).Find()
			res = append(res, SearchRespond{
				Mid: id.DocId,
			})
		}
		c.JSON(200, res)
	}
}
