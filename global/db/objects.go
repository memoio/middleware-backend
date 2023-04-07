package db

import "time"

type ObjectInfo struct {
	Address     string
	Name        string
	Size        int64
	Cid         string
	ModTime     time.Time
	CType       string
	UserDefined map[string]string
}
