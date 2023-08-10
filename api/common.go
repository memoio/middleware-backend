package api

import "time"

type ObjectOptions struct {
	Chain        int
	User         string
	Public       bool
	Key          []byte
	Size         int64
	MTime        time.Time
	DeleteMarker bool
	UserDefined  map[string]string
}
