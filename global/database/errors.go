package database

import (
	"fmt"
)

var (
	ErrAlreadyExist = fmt.Errorf("already exist")
	ErrNotExist     = fmt.Errorf("not exist")
)

type DataBase struct {
	message string
}

func (d DataBase) Error() string {
	return "DataBase Error: " + d.message
}
