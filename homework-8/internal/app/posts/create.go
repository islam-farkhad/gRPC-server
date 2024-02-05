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

// CreatePost handles the creation of a new post.
func (i *Implementation) CreatePost(ctx context.Context, req *pb.AddPostRequest) (*pb.Post, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "posts: CreatePost")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "posts/Implementation.CreatePost")))
	logger.Infof(ctx, "got request: %+v", req)

	if err := validateAddPostRequest(req); err != nil {
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, app.ValidationFailed, err)
	}

	postRepo := &repository.Post{
		Content: req.GetContent(),
		Likes:   req.GetLikes(),
	}

	postRepo, err := i.Repo.AddPost(ctx, postRepo)
	if err != nil {
		msg := "adding post error occurred"
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, msg, err)
	}

	return utils.ConvertRepoPostToProtoPost(postRepo), nil
}

func validateAddPostRequest(req *pb.AddPostRequest) error {
	if req.GetContent() == "" {
		return fmt.Errorf("%w: content cannot be empty", app.ErrValidationFail)
	}
	return nil
}
