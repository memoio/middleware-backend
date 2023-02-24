package gateway

import (
	"context"
	"fmt"
	"net/http"
)

const (
	newfunc     = "connect mefs error %s"
	makefunc    = "make bucket error %s"
	pricefunc   = "get price error %s"
	putfunc     = "put object error %s"
	getfunc     = "get object error %s"
	listfunc    = "list object error %s"
	deletefunc  = "delete object error %s"
	getinfofunc = "get object info error %s"
)

type StorageError struct {
	Storage string
	Message string
}

func (e StorageError) Error() string {
	return e.Storage + ":" + e.Message
}

func funcError(storage StorageType, efun string, err error) error {
	return StorageError{Storage: storage.String(), Message: fmt.Sprintf(efun, err)}
}

type NotImplemented struct {
	Message string
}

func (e NotImplemented) Error() string {
	return e.Message
}

type StorageNotSupport struct{}

func (e StorageNotSupport) Error() string {
	return "storage not support"
}

type AddressNull struct{}

func (e AddressNull) Error() string {
	return "address is nil"
}

type EthError struct {
	Message string
}

func (e EthError) Error() string {
	return e.Message
}

type BalanceNotEnough struct{}

func (e BalanceNotEnough) Error() string {
	return "balance not enough"
}

type APIError struct {
	Code           string
	Description    string
	HTTPStatusCode int
}

type APIErrorCode int

type errorCodeMap map[APIErrorCode]APIError

const (
	ErrNone APIErrorCode = iota
	ErrInternalError
	ErrNotImplemented
	ErrStorage
	ErrAddressNull
	ErrStorageNotSupport
	ErrBalanceNotEnough
)

func (e errorCodeMap) ToAPIErrWithErr(errCode APIErrorCode, err error) APIError {
	apiErr, ok := e[errCode]
	if !ok {
		apiErr = e[ErrInternalError]
	}
	if err != nil {
		apiErr.Description = fmt.Sprintf("%s (%s)", apiErr.Description, err)
	}
	return apiErr
}

func (e errorCodeMap) ToAPIErr(errCode APIErrorCode) APIError {
	return e.ToAPIErrWithErr(errCode, nil)
}

var ErrorCodes = errorCodeMap{
	ErrInternalError: {
		Code:           "InternalError",
		Description:    "We encountered an internal error, please try again.",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrNotImplemented: {
		Code:           "NotImplemented",
		Description:    "A header you provided implies functionality that is not implemented",
		HTTPStatusCode: http.StatusNotImplemented,
	},
	ErrStorage: {
		Code:           "Storage",
		Description:    "Error storing file",
		HTTPStatusCode: 516,
	},
	ErrAddressNull: {
		Code:           "Address",
		Description:    "Address Error",
		HTTPStatusCode: 517,
	},
	ErrStorageNotSupport: {
		Code:           "Storage",
		Description:    "Storage Error",
		HTTPStatusCode: 518,
	},
	ErrBalanceNotEnough: {
		Code:           "Balance",
		Description:    "Balance Error",
		HTTPStatusCode: 519,
	},
}

func ToAPIErrorCode(ctx context.Context, err error) (apiErr APIErrorCode) {
	if err == nil {
		return ErrNone
	}

	switch err.(type) {
	case NotImplemented:
		apiErr = ErrInternalError
	case StorageError:
		apiErr = ErrStorage
	case AddressNull:
		apiErr = ErrAddressNull
	case StorageNotSupport:
		apiErr = ErrStorageNotSupport
	case BalanceNotEnough:
		apiErr = ErrBalanceNotEnough
	}
	return apiErr
}
