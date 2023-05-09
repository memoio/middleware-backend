package gateway

import (
	"context"
	"io"
	"time"
)

type ObjectInfo struct {
	Address     string
	Name        string
	Size        int64
	Cid         string
	ModTime     time.Time
	CType       string
	UserDefined map[string]string
}

type ObjectOptions struct {
	Size         int64
	MTime        time.Time
	DeleteMarker bool
	UserDefined  map[string]string
}

type IGateway interface {
	PutObject(context.Context, string, string, io.Reader, ObjectOptions) (ObjectInfo, error)
	GetObject(context.Context, string, io.Writer, ObjectOptions) error
	ListObjects(context.Context, string) ([]ObjectInfo, error)
	GetObjectInfo(context.Context, string) (ObjectInfo, error)
}
