package siwb

import (
	"fmt"
)

type InvalidMessage struct{ string }

func (m *InvalidMessage) Error() string {
	return fmt.Sprintf("Invalid Message: %s", m.string)
}
