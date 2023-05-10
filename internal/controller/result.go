package controller

import "time"

type PutObjectResult struct {
	Mid string
}

type GetObjectResult struct {
	Name  string
	Size  int64
	CType string
}

type ListObjectsResult struct {
	Address string
	Storage string
	Objects []ObjectInfoResult
}

type ObjectInfoResult struct {
	Name        string
	Size        int64
	Mid         string
	ModTime     time.Time
	UserDefined map[string]string
}
