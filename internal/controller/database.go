package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/logs"
)

func (c *Controller) ListObjects(ctx context.Context, chain int, address string) (ListObjectsResult, error) {
	result := ListObjectsResult{}

	loi, err := database.List(chain, address, c.storageType)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Storage = c.storageType.String()

	for _, oi := range loi {
		userdefine := make(map[string]string)
		err = json.Unmarshal([]byte(oi.UserDefine), &userdefine)
		if err != nil {
			lerr := logs.ControllerError{Message: fmt.Sprint("unmarshal userdefine error, ", err)}
			logger.Error(lerr)
			return result, lerr
		}

		result.Objects = append(result.Objects, ObjectInfoResult{
			Name:    oi.Name,
			Size:    oi.Size,
			Mid:     oi.Mid,
			ModTime: oi.ModTime,
			Public:  oi.Public,
			// UserDefined: userdefine,
		})
	}

	return result, nil
}

func (c *Controller) GetObjectInfo(ctx context.Context, chain int, mid string) (map[string]database.FileInfo, error) {
	return database.Get(chain, mid, c.storageType)
}
