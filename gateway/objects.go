package gateway

import (
	"time"
)

type ObjectOptions struct {
	Size         int64
	MTime        time.Time
	DeleteMarker bool
	UserDefined  map[string]string
}


