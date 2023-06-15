package server

import (
	"github.com/memoio/backend/internal/contract"
	"github.com/memoio/backend/internal/gateway"
)

type api struct {
	gateway  gateway.IGateway
	contract *contract.Contract
}
