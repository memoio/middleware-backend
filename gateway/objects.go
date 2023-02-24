package gateway

import "time"

type ObjectOptions struct {
	PayType      string
	MTime        time.Time
	DeleteMarker bool
	UserDefined  map[string]string
}

type ObjectInfo struct {
	Address     string
	Name        string
	Size        int64
	Cid         string
	ModTime     time.Time
	CType       string
	UserDefined map[string]string
}

type ListObjectsInfo struct {
	Objects []ObjectInfo
}

type StorageInfo struct {
	Used      string
	Available string
	Free      string
	Files     string
}
