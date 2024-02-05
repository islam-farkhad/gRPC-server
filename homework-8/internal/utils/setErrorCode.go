package utils

import (
	"errors"
	"homework-8/internal/app"
	"homework-8/internal/pkg/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetErrorCode performs error checking ands sets appropriate error code.
func SetErrorCode(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrObjectNotFound) {
		err = status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, app.ErrValidationFail) {
		err = status.Error(codes.InvalidArgument, err.Error())
	}
	// Maybe some other checks

	return status.Error(codes.Internal, err.Error())
}
