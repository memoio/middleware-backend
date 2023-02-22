package api

import (
	"context"
	"io"

	"github.com/memoio/backend/gateway"
)

type ObjectLayer interface {
	Putobject(context.Context, string, string, io.Reader) (gateway.ObjectInfo, error)
	GetobjectNInfo(context.Context, string, string) (io.Reader, error)
}
