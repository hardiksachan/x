// see also: https://middlemost.com/failure-is-your-domain/
// see also: https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html

package xerrors

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/Logistics-Coordinators/x/xlog"
)

type (
	Op      string
	Message string
	Code    uint8
)

const (
	Other Code = iota
	Internal
	Invalid
	NotFound
	Exists
	Expired
)

func (c *Code) String() string {
	switch *c {
	case Other:
		return "other error"
	case Internal:
		return "internal error"
	case Invalid:
		return "invalid error"
	case NotFound:
		return "item not found"
	case Exists:
		return "item already exists"
	case Expired:
		return "item has expired"
	}
	return "unknown error code"
}

type Error struct {
	// Machine-readable error code.
	Code Code

	// Human-readable message.
	Message Message

	// Logical operation and nested error.
	Op  Op
	Err error
}

// Error returns the string representation of the error message.
func (e *Error) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		_, _ = fmt.Fprintf(&buf, "%s: ", e.Op)
	}

	if e.Code != Other {
		_, _ = fmt.Fprintf(&buf, "<%s> ", e.Code.String())
	}

	if e.Message != "" {
		_, _ = fmt.Fprintf(&buf, "%s", e.Message)
	}

	if e.Err != nil {
		_, _ = fmt.Fprintf(&buf, "\n\t\t%s", e.Err.Error())
	}

	return buf.String()
}

// ErrorCode returns the code of the root error, if available. Otherwise returns Internal.
func ErrorCode(err error) Code {
	if err == nil {
		return Other
	} else if e, ok := err.(*Error); ok && e.Code != Other {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(e.Err)
	}
	return Internal
}

// ErrorMessage returns the human-readable message of the error, if available.
// Otherwise, returns a generic error message.
func ErrorMessage(err error) Message {
	if err == nil {
		return ""
	} else if e, ok := err.(*Error); ok && e.Message != "" {
		return e.Message
	} else if ok && e.Err != nil {
		return ErrorMessage(e.Err)
	}
	return "An internal error has occurred. Please contact technical support."
}

func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Message:
			e.Message = arg
		case Code:
			e.Code = arg
		case *Error:
			argCopy := *arg
			e.Err = &argCopy
		case error:
			e.Err = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			xlog.Errorf("errors.E: bad call from %s:%d: %v", file, line, args)
			panic(fmt.Sprintf("unknown type %T, value %v in error call", arg, arg))
		}
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind, or file name twice.
	if prev.Code == e.Code {
		prev.Code = Other
	}
	// If this error has Kind unset or Other, pull up the inner one.
	if e.Code == Other {
		e.Code = prev.Code
		prev.Code = Other
	}
	return e
}
