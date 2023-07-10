package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/database"
	"github.com/memoio/backend/internal/logs"
)

func storeFileInfo() error {
	return nil
}

func (c *Controller) getObjectInfo(ctx context.Context, address, mid string, st api.StorageType) (database.FileInfo, error) {
	result := database.FileInfo{}
	oi, err := c.database.GetObjectInfo(ctx, address, mid, st)
	if err != nil {
		return result, err
	}

	fi := oi.(database.FileInfo)
	return fi, nil
}

func (c *Controller) getObjectInfoById(ctx context.Context, id int) (database.FileInfo, error) {
	result := database.FileInfo{}
	oi, err := c.database.GetObjectInfoById(ctx, id)
	if err != nil {
		return result, err
	}

	fi := oi.(database.FileInfo)
	return fi, nil
}

func storeFlowSize() error {
	return nil
}

func (c *Controller) listobjects(ctx context.Context, address string) (ListObjectsResult, error) {
	result := ListObjectsResult{}

	loi, err := c.database.ListObjects(ctx, address, c.st)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Storage = c.st.String()

	for _, ioi := range loi {
		userdefine := make(map[string]string)
		oi := ioi.(database.FileInfo)
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

func (c *Controller) deleteObject(ctx context.Context, address, mid string) error {
	return logs.NotImplemented{}
}

func (c *Controller) getCacheStorageInfo(ctx context.Context, address string) (int64, error) {
	return 0, logs.NotImplemented{}
}
