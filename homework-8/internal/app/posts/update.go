package posts

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

// UpdatePost manages the update of an existing post.
func (i *Implementation) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "posts: UpdatePost")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "posts/Implementation.UpdatePost")))
	logger.Infof(ctx, "got request: %+v", req)

	if err := validateUpdatePostRequest(req); err != nil {
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, app.ValidationFailed, err)
	}
	postRepo := &repository.Post{
		ID:      req.GetId(),
		Content: req.GetContent(),
		Likes:   req.GetLikes(),
	}

	postRepo, err := i.Repo.UpdatePost(ctx, postRepo)
	if err != nil {
		msg := "could not update post"
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, msg, err)
	}

	return utils.ConvertRepoPostToProtoPost(postRepo), nil
}

func validateUpdatePostRequest(req *pb.UpdatePostRequest) error {
	if req.Id < 1 {
		return fmt.Errorf("%w: id should be positive number: %d", app.ErrValidationFail, req.GetId())
	}

	if req.Content == "" {
		return fmt.Errorf("%w: updated content must be provided", app.ErrValidationFail)
	}

	return nil
}
