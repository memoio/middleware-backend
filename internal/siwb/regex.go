package siwb

import (
	"fmt"
	"regexp"
)

const _SIWE_DOMAIN = "(?P<domain>([^/?#]+)) wants you to sign in with your Bitcoin account:\\n"
const _SIWE_ADDRESS = "(?P<address>[a-zA-Z0-9]+)\\n\\n"

// const _SIWE_STATEMENT = "((?P<statement>[^\\n]+)\\n)?\\n"
const _RFC3986 = "(([^ :/?#]+):)?(//([^ /?#]*))?([^ ?#]*)(\\?([^ #]*))?(#(.*))?"

var _SIWE_URI_LINE = fmt.Sprintf("URI: (?P<uri>%s?)\\n", _RFC3986)

const _SIWE_VERSION = "Version: (?P<version>1)\\n"
const _SIWE_NONCE = "Nonce: (?P<nonce>[a-zA-Z0-9]{8,})\\n"
const _SIWE_DATETIME = "([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))"

var _SIWE_ISSUED_AT = fmt.Sprintf("Issued At: (?P<issuedAt>%s)", _SIWE_DATETIME)

var _SIWE_MESSAGE = regexp.MustCompile(fmt.Sprintf("^%s%s%s%s%s%s$",
	_SIWE_DOMAIN,
	_SIWE_ADDRESS,
	// _SIWE_STATEMENT,
	_SIWE_URI_LINE,
	_SIWE_VERSION,
	_SIWE_NONCE,
	_SIWE_ISSUED_AT))
