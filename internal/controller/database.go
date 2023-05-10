package controller

import (
	"context"

	"github.com/memoio/backend/internal/database"
)

func (c *Controller) ListObjects(ctx context.Context, address string) (ListObjectsResult, error) {
	result := ListObjectsResult{}

	loi, err := database.List(address, c.storageType)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Storage = c.storageType.String()

	for _, oi := range loi {
		result.Objects = append(result.Objects, ObjectInfoResult{
			Name: oi.Name,
			Size: oi.Size,
			Mid:  oi.Mid,
		})
	}

	return result, nil
}
