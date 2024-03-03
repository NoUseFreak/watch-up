package watchup

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func codeFromError(err error) int {
	urlErr, ok := err.(*url.Error)
	if !ok {
		return E_UNKNOWN
	}
	opErr, ok := urlErr.Err.(*net.OpError)
	if !ok {
		return E_UNKNOWN
	}

	switch opErr.Err.(type) {
	case *net.DNSError:
		return E_DNS
	case *net.AddrError:
		return E_ADDR
	case *net.UnknownNetworkError:
		return E_NET
	case *net.InvalidAddrError:
		return E_INVALID
	case *net.ParseError:
		return E_PARSE
	case *net.OpError:
		return E_OP
	case *os.SyscallError:
		if strings.Contains(opErr.Err.Error(), "connection refused") {
			return E_REFUSED
		}
		return E_SYSCALL
	default:
		logrus.Errorf("Unhandled err: %s", opErr.Err.Error())
		return E_UNKNOWN
	}
}

const (
	E_UNKNOWN = 0
	E_DNS     = 1
	E_ADDR    = 2
	E_NET     = 3
	E_INVALID = 4
	E_PARSE   = 5
	E_OP      = 6
	E_SYSCALL = 7
	E_REFUSED = 8
)

// CodeToText returns the text representation of the error code.
func CodeToText(code int) string {
	text := http.StatusText(code)
	if text != "" {
		return text
	}

	switch code {
	case E_UNKNOWN:
		return "Unknown error"
	case E_DNS:
		return "DNS error"
	case E_ADDR:
		return "Address error"
	case E_NET:
		return "Network error"
	case E_INVALID:
		return "Invalid address error"
	case E_PARSE:
		return "Parse error"
	case E_OP:
		return "Operation error"
	case E_SYSCALL:
		return "System call error"
	case E_REFUSED:
		return "Connection refused"
	default:
		logrus.Errorf("Unknown error code: %d", code)
		return "Unknown error"
	}
}
