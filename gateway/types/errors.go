package types

type Logs struct {
	method  string
	message string
}

func New(method string) Logs {
	return Logs{method: method}
}
func (l Logs) DealError(e error) error {
	l.message = e.Error()
	return l
}

func (l Logs) Error() string {
	return l.method + " : " + l.message
}
