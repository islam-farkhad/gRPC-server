package utils

import (
	"context"
	"fmt"
	"homework-8/internal/pkg/logger"
)

// ReportError helper function - logs and returns error. Expected such  format: `message %w`
func ReportError(ctx context.Context, msg string, err error) error {
	logger.Errorf(ctx, "%s: %+v", msg, err)
	return fmt.Errorf("%s: %w", msg, err)
}
