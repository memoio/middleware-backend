package server

import "time"

type BalanceResponse struct {
	Address string
	Balance string
}

type ListObjectsResponse struct {
	Address string
	Storage string
	Object  []ObjectResponse
}

type ObjectResponse struct {
	Name        string
	Size        int64
	Cid         string
	ModTime     time.Time
	UserDefined map[string]string
}
