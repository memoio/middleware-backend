package logs

import "fmt"

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

func (g *GenericError) SetMethod(m string) {
	g.method = m
}

func (g *GenericError) DealError(e error) error {
	g.message = e.Error()
	return g
}

func (g *GenericError) Error() string {
	return g.method + " : " + g.message
}

func (g *GenericError) ToResponse() ErrResponse {
	return ErrResponse{Package: g.pkg, Err: g}
}
