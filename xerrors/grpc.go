package xerrors

import (
	"github.com/Logistics-Coordinators/x/xlog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c Code) grpcCode() codes.Code {
	switch c {
	case Other:
		return codes.Unknown
	case Internal:
		return codes.Internal
	case Invalid:
		return codes.InvalidArgument
	case NotFound:
		return codes.NotFound
	case Exists:
		return codes.AlreadyExists
	}
	return codes.Unknown
}

func GrpcError(err error) error {
	code := ErrorCode(err)

	xlog.ErrorString(err.Error())

	grpcCode := code.grpcCode()
	message := ErrorMessage(err)

	return status.Errorf(grpcCode, string(message))
}
