package posts

import (
	"context"
	"fmt"
	"homework-8/internal/app"
	"homework-8/internal/pkg/logger"
	"homework-8/internal/utils"
	pb "homework-8/pkg/posts"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// DeletePost manages the deletion of a post by its ID.
func (i *Implementation) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*empty.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "posts: DeletePost")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "posts/Implementation.DeletePost")))
	logger.Infof(ctx, "got request: %+v", req)

	if err := validateDeletePostRequest(req); err != nil {
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, app.ValidationFailed, err)
	}

	if err := i.Repo.DeletePostByID(ctx, req.GetId()); err != nil {
		msg := "deleting post error occurred"
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, msg, err)
	}

	return &empty.Empty{}, nil
}

func validateDeletePostRequest(req *pb.DeletePostRequest) error {
	if req.Id < 1 {
		return fmt.Errorf("%w: id should be positive number: %d", app.ErrValidationFail, req.GetId())
	}
	return nil
}
