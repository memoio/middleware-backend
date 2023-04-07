package gateway

import (
	"time"

	"github.com/memoio/backend/global/db"
)

type ObjectOptions struct {
	Size         int64
	MTime        time.Time
	DeleteMarker bool
	UserDefined  map[string]string
}

type ObjectInfo db.ObjectInfo


