package api

import (
	"context"
	"io"
)

type IGateway interface {
	PutObject(context.Context, string, string, io.Reader, ObjectOptions) (ObjectInfo, error)
	GetObject(context.Context, string, io.Writer, ObjectOptions) error
	DeleteObject(context.Context, string, string) error
	// ListObjects(context.Context, string) ([]ObjectInfo, error)
	// GetObjectInfo(context.Context, string) (ObjectInfo, error)
}
