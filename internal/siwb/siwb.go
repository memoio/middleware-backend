package siwb

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func buildAuthority(uri *url.URL) string {
	authority := uri.Host
	if uri.User != nil {
		authority = fmt.Sprintf("%s@%s", uri.User.String(), authority)
	}
	return authority
}

func validateDomain(domain *string) (bool, error) {
	if isEmpty(domain) {
		return false, &InvalidMessage{"`domain` must not be empty"}
	}

	validateDomain, err := url.Parse(fmt.Sprintf("https://%s", *domain))
	if err != nil {
		return false, &InvalidMessage{"Invalid format for field `domain`"}
	}

	authority := buildAuthority(validateDomain)
	if authority != *domain {
		return false, &InvalidMessage{"Invalid format for field `domain`"}
	}

	return true, nil
}

func validateURI(uri *string) (*url.URL, error) {
	if isEmpty(uri) {
		return nil, &InvalidMessage{"`uri` must not be empty"}
	}

	validateURI, err := url.Parse(*uri)
	if err != nil {
		return nil, &InvalidMessage{"Invalid format for field `uri`"}
	}

	return validateURI, nil
}

// InitMessage creates a Message object with the provided parameters
func InitMessage(domain, address, uri, nonce string, options map[string]interface{}) (*Message, error) {
	if ok, err := validateDomain(&domain); !ok {
		return nil, err
	}

	if isEmpty(&address) {
		return nil, &InvalidMessage{"`address` must not be empty"}
	}

	validateURI, err := validateURI(&uri)
	if err != nil {
		return nil, err
	}

	if isEmpty(&nonce) {
		return nil, &InvalidMessage{"`nonce` must not be empty"}
	}

	// var statement *string
	// if val, ok := options["statement"]; ok {
	// 	value := val.(string)
	// 	statement = &value
	// }

	var issuedAt string
	timestamp, err := parseTimestamp(options, "issuedAt")
	if err != nil {
		return nil, err
	}

	if timestamp != nil {
		issuedAt = *timestamp
	} else {
		issuedAt = time.Now().UTC().Format(time.RFC3339)
	}

	return &Message{
		domain:  domain,
		address: address,
		uri:     *validateURI,
		version: "1",

		// statement: statement,
		nonce: nonce,

		issuedAt: issuedAt,
	}, nil
}

func parseMessage(message string) (map[string]interface{}, error) {
	match := _SIWE_MESSAGE.FindStringSubmatch(message)

	if match == nil {
		return nil, &InvalidMessage{"Message could not be parsed"}
	}

	result := make(map[string]interface{})
	for i, name := range _SIWE_MESSAGE.SubexpNames() {
		if i != 0 && name != "" && match[i] != "" {
			result[name] = match[i]
		}
	}

	if _, ok := result["domain"]; !ok {
		return nil, &InvalidMessage{"`domain` must not be empty"}
	}
	domain := result["domain"].(string)
	if ok, err := validateDomain(&domain); !ok {
		return nil, err
	}

	if _, ok := result["uri"]; !ok {
		return nil, &InvalidMessage{"`domain` must not be empty"}
	}
	uri := result["uri"].(string)
	if _, err := validateURI(&uri); err != nil {
		return nil, err
	}

	return result, nil
}

// ParseMessage returns a Message object by parsing an EIP-4361 formatted string
func ParseMessage(message string) (*Message, error) {
	result, err := parseMessage(message)
	if err != nil {
		return nil, err
	}

	parsed, err := InitMessage(
		result["domain"].(string),
		result["address"].(string),
		result["uri"].(string),
		result["nonce"].(string),
		result,
	)

	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func (m *Message) MessageHash() (hash []byte, err error) {
	var buf bytes.Buffer
	if err = wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n"); err != nil {
		return nil, err
	}
	if err = wire.WriteVarString(&buf, 0, m.String()); err != nil {
		return nil, err
	}

	hash = chainhash.DoubleHashB(buf.Bytes())
	return hash, nil
}

func (m *Message) prepareMessage() string {
	greeting := fmt.Sprintf("%s wants you to sign in with your Bitcoin account:", m.domain)
	headerArr := []string{greeting, m.address}

	// if isEmpty(m.statement) {
	// 	headerArr = append(headerArr, "\n")
	// } else {
	// 	headerArr = append(headerArr, fmt.Sprintf("\n%s\n", *m.statement))
	// }

	headerArr = append(headerArr, "")
	header := strings.Join(headerArr, "\n")

	uri := fmt.Sprintf("URI: %s", m.uri.String())
	version := fmt.Sprintf("Version: %s", m.version)
	nonce := fmt.Sprintf("Nonce: %s", m.nonce)
	issuedAt := fmt.Sprintf("Issued At: %s", m.issuedAt)

	bodyArr := []string{uri, version, nonce, issuedAt}
	body := strings.Join(bodyArr, "\n")

	return strings.Join([]string{header, body}, "\n")
}

func (m *Message) String() string {
	return m.prepareMessage()
}
