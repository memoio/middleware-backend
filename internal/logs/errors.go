package logs

import (
	"fmt"
	"net/http"
)

type StorageError struct {
	Storage string
	Message string
}

func (e StorageError) Error() string {
	return e.Storage + ":" + e.Message
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

type AddressError struct {
	Message string
}

func (e AddressError) Error() string {
	return e.Message
}

type AuthenticationFailed struct {
	Message string
}

func (e AuthenticationFailed) Error() string {
	return e.Message
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

type ServerError struct {
	Message string
}

func (e ServerError) Error() string {
	return e.Message
}

type GatewayError struct {
	Message string
}

func (e GatewayError) Error() string {
	return e.Message
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
	ErrAddressError
	ErrStorageNotSupport
	ErrAuthenticationFailed
	ErrBalanceNotEnough
	ErrEthError
	ErrServerError
	ErrGatewayError
)

func (e errorCodeMap) ToAPIErrWithErr(errCode APIErrorCode, err error) APIError {
	apiErr, ok := e[errCode]
	if !ok {
		apiErr = e[ErrInternalError]
	}
	if err != nil {
		apiErr.Description = fmt.Sprintf("%s (%s)", apiErr.Description, err.Error())
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
	ErrAddressError: {
		Code:           "Address",
		Description:    "Address Error",
		HTTPStatusCode: 517,
	},
	ErrStorageNotSupport: {
		Code:           "Storage",
		Description:    "Storage Error",
		HTTPStatusCode: 518,
	},
	ErrAuthenticationFailed: {
		Code:           "Authentication",
		Description:    "Authentication Failed",
		HTTPStatusCode: 401,
	},
	ErrBalanceNotEnough: {
		Code:           "Balance",
		Description:    "Balance Error",
		HTTPStatusCode: 519,
	},
	ErrEthError: {
		Code:           "Eth",
		Description:    "Eth Error",
		HTTPStatusCode: 520,
	},
	ErrServerError: {
		Code:           "ServerError",
		Description:    "Server Error",
		HTTPStatusCode: 521,
	},
	ErrGatewayError: {
		Code:           "GatewayError",
		Description:    "Gateway Error",
		HTTPStatusCode: 522,
	},
}

func ToAPIErrorCode(err error) APIError {
	if err == nil {
		return ErrorCodes.ToAPIErr(ErrNone)
	}
	var apiErr APIErrorCode

	switch err.(type) {
	case NotImplemented:
		apiErr = ErrNotImplemented
	case StorageError:
		apiErr = ErrStorage
	case AddressError:
		apiErr = ErrAddressError
	case StorageNotSupport:
		apiErr = ErrStorageNotSupport
	case AuthenticationFailed:
		apiErr = ErrAuthenticationFailed
	case BalanceNotEnough:
		apiErr = ErrBalanceNotEnough
	case EthError:
		apiErr = ErrEthError
	default:
		apiErr = ErrInternalError
	}
	return ErrorCodes.ToAPIErrWithErr(apiErr, err)
}

var (
	ErrAlreadyExist = fmt.Errorf("already exist")
	ErrNotExist     = fmt.Errorf("not exist")
)

type ErrResponse struct {
	Package string
	Err     error
}

type GenericError struct {
	pkg     string
	method  string
	message string
}

func SetPkg(pkg string) GenericError {
	return GenericError{pkg: pkg}
}
