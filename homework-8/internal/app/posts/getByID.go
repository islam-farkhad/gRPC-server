package posts

import (
	"context"
	"fmt"
	"homework-8/internal/app"
	"homework-8/internal/pkg/logger"
	"homework-8/internal/utils"
	pb "homework-8/pkg/posts"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// GetPostByID retrieves a post by its ID and its associated comments.
func (i *Implementation) GetPostByID(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostByIDResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "posts: GetPostByID")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "posts/Implementation.GetPostByID")))
	logger.Infof(ctx, "got request: %+v", req)

	if err := validateGetPostRequest(req); err != nil {
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, app.ValidationFailed, err)
	}

	postRepo, err := i.Repo.GetPostByID(ctx, req.GetId())
	if err != nil {
		msg := "getting post by id error occurred"
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, msg, err)
	}

	comments, err := i.Repo.GetCommentsByPostID(ctx, req.GetId())
	if err != nil {
		msg := "getting comments by post id error occurred"
		err = utils.SetErrorCode(err)
		return nil, utils.ReportError(ctx, msg, err)
	}

	var protoComments []*pb.Comment
	for _, comment := range comments {
		protoComments = append(protoComments, utils.ConvertRepoCommentToProtoComment(comment))
	}

	return &pb.GetPostByIDResponse{
		Post:     utils.ConvertRepoPostToProtoPost(postRepo),
		Comments: protoComments,
	}, nil
}

func validateGetPostRequest(req *pb.GetPostRequest) error {
	if req.Id < 1 {
		return fmt.Errorf("%w: id should be positive number: %d", app.ErrValidationFail, req.GetId())
	}
	return nil
}
