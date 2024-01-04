package siwb

import (
	"net/url"
)

type Message struct {
	domain  string
	address string
	uri     url.URL
	version string

	// statement *string
	nonce string

	issuedAt string
}

func (m *Message) GetDomain() string {
	return m.domain
}

func (m *Message) GetAddress() string {
	return m.address
}

func (m *Message) GetURI() url.URL {
	return m.uri
}

func (m *Message) GetVersion() string {
	return m.version
}

// func (m *Message) GetStatement() *string {
// 	if m.statement != nil {
// 		ret := *m.statement
// 		return &ret
// 	}
// 	return nil
// }

func (m *Message) GetNonce() string {
	return m.nonce
}

func (m *Message) GetIssuedAt() string {
	return m.issuedAt
}
