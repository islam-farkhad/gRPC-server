package comments

import (
	"context"
	"fmt"
	"homework-8/internal/app"
	"homework-8/internal/pkg/logger"
	"homework-8/internal/pkg/repository"
	"homework-8/internal/utils"
	pb "homework-8/pkg/posts"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// CreateComment handles the creation of a new comment.
func (i *Implementation) CreateComment(ctx context.Context, req *pb.AddCommentRequest) (*pb.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "comments: CreateComment")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "CreateComment")))
	logger.Infof(ctx, "got request: %+v", req)

	if err := validateAddCommentRequest(req); err != nil {
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, app.ValidationFailed, err)
	}

	commentRepo := &repository.Comment{
		PostID:  req.GetPostId(),
		Content: req.GetContent(),
	}
	commentRepo, err := i.Repo.AddComment(ctx, commentRepo)
	if err != nil {
		msg := "error adding comment occurred"
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, msg, err)
	}

	return utils.ConvertRepoCommentToProtoComment(*commentRepo), nil
}

func validateAddCommentRequest(req *pb.AddCommentRequest) error {
	if req.GetPostId() < 1 {
		return fmt.Errorf("%w: id should be positive number: %d", app.ErrValidationFail, req.GetPostId())
	}

	if req.GetContent() == "" {
		return fmt.Errorf("%w: content cannot be empty", app.ErrValidationFail)
	}
	return nil
}
