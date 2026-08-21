package error

import (
	"errors"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"
	"testing"

	"github.com/googleapis/gax-go/v2/apierror"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestError_ErrWithStackTrace(t *testing.T) {
	err := errors.New("this is error")

	tt := []struct {
		name    string
		err     error
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Should return success convert error to pkgerrors.WithStack",
			err:  err,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				type stackTracer interface {
					StackTrace() pkgerrors.StackTrace
				}

				_, ok := err.(stackTracer)
				return assert.EqualError(t, i[0].(error), err.Error()) && assert.True(t, ok)
			},
		},
		{
			name: "Should return success if error already with pkgerrors type",
			err:  pkgerrors.WithStack(err),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				type stackTracer interface {
					StackTrace() pkgerrors.StackTrace
				}

				_, ok := err.(stackTracer)
				return assert.EqualError(t, i[0].(error), err.Error()) && assert.True(t, ok)
			},
		},
		{
			name: "Should return nil",
			err:  nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, i[0]) && assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.wantErr(t, ErrWithStackTrace(tc.err), tc.err)
		})
	}
}

func TestError_ErrCause(t *testing.T) {
	tt := []struct {
		name    string
		err     func() error
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Should return error with customerror.CustomError type",
			err: func() error {
				err := customerror.ErrBadRequest
				return ErrWithStackTrace(err)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, customerror.ErrBadRequest, err.Error())
			},
		},
		{
			name: "Should return customerror.ErrInternalServer error",
			err: func() error {
				err := errors.New("this is error")
				return ErrWithStackTrace(err)
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.EqualError(t, customerror.ErrInternalServer, err.Error())
			},
		},
		{
			name: "Should return nil",
			err: func() error {
				return nil
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.wantErr(t, ErrCause(tc.err()))
		})
	}
}

func TestError_TranslateAPIErrorToInternalError(t *testing.T) {
	tt := []struct {
		name        string
		err         func() error
		internalErr error
		expected    func() error
	}{
		{
			name: "Should return error with *apierror.APIError != NotFound",
			err: func() error {
				theErr := status.New(codes.Unknown, "unknown").Err()
				apiErr, _ := apierror.ParseError(theErr, true)
				return apiErr
			},
			internalErr: customerror.NewClientError(customerror.CodeBadRequest, "Invalid Request, Please Check Your Request Body."),
			expected: func() error {
				theErr := status.New(codes.Unknown, "unknown").Err()
				apiErr, _ := apierror.ParseError(theErr, true)
				return apiErr
			},
		},
		{
			name: "Should return error with *apierror.APIError = NotFound",
			err: func() error {
				theErr := status.New(codes.NotFound, "not found").Err()
				apiErr, _ := apierror.ParseError(theErr, true)
				return apiErr
			},
			internalErr: customerror.NewClientError(customerror.CodeBadRequest, "Invalid Request, Please Check Your Request Body."),
			expected: func() error {
				return customerror.NewClientError(customerror.CodeBadRequest, "Invalid Request, Please Check Your Request Body.")
			},
		},
		{
			name: "Should return error default",
			err: func() error {
				return customerror.NewClientError(customerror.CodeBadRequest, "default error.")
			},
			internalErr: customerror.NewClientError(customerror.CodeBadRequest, "Invalid Request, Please Check Your Request Body."),
			expected: func() error {
				return customerror.NewClientError(customerror.CodeBadRequest, "default error.")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected(), TranslateAPIErrorToInternalError(tc.err(), tc.internalErr))
		})
	}
}
