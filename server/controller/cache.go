package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/memoio/backend/api"
	"github.com/memoio/backend/internal/logs"
)

func (c *Controller) storeFileInfo(ctx context.Context, fi api.FileInfo) error {
	return c.database.PutObject(ctx, fi)
}

func (c *Controller) getObjectInfo(ctx context.Context, address, mid string, st api.StorageType) (api.FileInfo, error) {
	result := api.FileInfo{}
	oi, err := c.database.GetObjectInfo(ctx, address, mid, st)
	if err != nil {
		return result, err
	}

	fi := oi.(api.FileInfo)
	if fi == result {
		lerr := logs.DataBaseError{Message: "file not exist"}
		logger.Error(lerr)
		return result, lerr
	}
	return fi, nil
}

func (c *Controller) getObjectInfoById(ctx context.Context, id int) (api.FileInfo, error) {
	result := api.FileInfo{}
	oi, err := c.database.GetObjectInfoById(ctx, id)
	if err != nil {
		return result, err
	}

	fi := oi.(api.FileInfo)
	if fi == result {
		lerr := logs.DataBaseError{Message: "file not exist"}
		logger.Error(lerr)
		return result, lerr
	}
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
		oi := ioi.(api.FileInfo)
		err = json.Unmarshal([]byte(oi.UserDefine), &userdefine)
		if err != nil {
			lerr := logs.ControllerError{Message: fmt.Sprint("unmarshal userdefine error, ", err)}
			logger.Error(lerr)
			return result, lerr
		}

		result.Objects = append(result.Objects, ObjectInfoResult{
			ID:      oi.ID,
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

func (c *Controller) deleteObject(ctx context.Context, id int) error {
	return c.database.DeleteObject(ctx, id)
}

func (c *Controller) getCacheStorageInfo(ctx context.Context, address string) (int64, error) {
	return 0, logs.NotImplemented{}
}
