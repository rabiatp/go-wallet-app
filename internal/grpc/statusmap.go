package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/profiling/service"
	"google.golang.org/grpc/status"
)

func toStatus(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidAmount):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrWalletNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, service.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, service.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrFailedPrecondition):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, service.ErrAborted):
		return status.Error(codes.Aborted, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
