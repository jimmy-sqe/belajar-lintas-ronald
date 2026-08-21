package error

import (
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"

	"github.com/googleapis/gax-go/v2/apierror"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

func ErrWithStackTrace(err error) error {
	if err == nil {
		return nil
	}

	type stackTracer interface {
		StackTrace() pkgerrors.StackTrace
	}

	_, ok := err.(stackTracer)
	if !ok {
		return pkgerrors.WithStack(err)
	}

	return err
}

func ErrCause(err error) error {
	if err == nil {
		return nil
	}

	switch causeErr := pkgerrors.Cause(err).(type) {
	case customerror.CustomError:
		return causeErr
	default:
		return customerror.ErrInternalServer
	}
}

func TranslateAPIErrorToInternalError(err, internalErr error) error {
	switch causeErr := pkgerrors.Cause(err).(type) {
	case *apierror.APIError:
		code := causeErr.GRPCStatus().Code()
		if code == codes.NotFound {
			return internalErr
		}
		return err
	default:
		return err
	}
}
