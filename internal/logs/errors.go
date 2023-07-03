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

type ContractError struct {
	Message string
}

func (e ContractError) Error() string {
	return e.Message
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

type ConfigError struct {
	Message string
}

func (e ConfigError) Error() string {
	return e.Message
}

type DataBaseError struct {
	Message string
}

func (e DataBaseError) Error() string {
	return e.Message
}

type ControllerError struct {
	Message string
}

func (e ControllerError) Error() string {
	return e.Message
}

type NoPermission struct {
	Message string
}

func (e NoPermission) Error() string {
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
	ErrContractError
	ErrEthError
	ErrServerError
	ErrGatewayError
	ErrConfigError
	ErrDataBaseError
	ErrControllerError
	ErrNoPermission
)

func (e errorCodeMap) ToAPIErrWithErr(errCode APIErrorCode, err error) APIError {
	apiErr, ok := e[errCode]
	if !ok {
		apiErr = e[ErrAddressError]
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
	ErrContractError: {
		Code:           "contract",
		Description:    "contract Error",
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
	ErrConfigError: {
		Code:           "ConfigError",
		Description:    "Config Error",
		HTTPStatusCode: 523,
	},
	ErrDataBaseError: {
		Code:           "DataBaseError",
		Description:    "DataBase Error",
		HTTPStatusCode: 524,
	},
	ErrControllerError: {
		Code:           "ControllerError",
		Description:    "Controller Error",
		HTTPStatusCode: 525,
	},
	ErrNoPermission: {
		Code:           "Permission",
		Description:    "You don't have access to the resource",
		HTTPStatusCode: 526,
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
	case ControllerError:
		apiErr = ErrControllerError
	case AddressError:
		apiErr = ErrAddressError
	case StorageNotSupport:
		apiErr = ErrStorageNotSupport
	case AuthenticationFailed:
		apiErr = ErrAuthenticationFailed
	case ContractError:
		apiErr = ErrContractError
	case ServerError:
		apiErr = ErrServerError
	case EthError:
		apiErr = ErrEthError
	case GatewayError:
		apiErr = ErrGatewayError
	case ConfigError:
		apiErr = ErrConfigError
	case DataBaseError:
		apiErr = ErrDataBaseError
	case NoPermission:
		apiErr = ErrNoPermission
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
