package main

import (
	"context"
	"fmt"
	"homework-8/internal/config"
	"homework-8/internal/pkg/logger"
	pb "homework-8/pkg/posts"
	"log"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// All service supported methods.
const (
	method = CreatePost // Choose method here

	CreatePost  = "CreatePost"
	GetPostByID = "GetPostByID"
	UpdatePost  = "UpdatePost"
	DeletePost  = "DeletePost"

	CreateComment = "CreateComment"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("error creating zap logger")
	}
	logger.SetGlobal(
		zapLogger.With(zap.String("component", "client")),
	)

	addr := config.GetConfigs().GetPort()
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error trying to make connection %w", err)
	}

	logger.Infof(ctx, "Connection with %q is set", addr)
	postClient := pb.NewPostServiceClient(conn)
	commentClient := pb.NewCommentServiceClient(conn)

	switch method {
	case CreatePost:
		return remoteCallCreatePost(ctx, postClient)

	case GetPostByID:
		return remoteCallGetPostByID(ctx, postClient)

	case UpdatePost:
		return remoteCallUpdatePost(ctx, postClient)

	case DeletePost:
		return remoteCallDeletePost(ctx, postClient)

	case CreateComment:
		return remoteCallCreateComment(ctx, commentClient)
	default:
		msg := fmt.Sprintf("unknown method %s (use: %s, %s, %s, %s, %s)", method, CreatePost, GetPostByID, UpdatePost, DeletePost, CreateComment)
		logger.Errorf(ctx, msg)
		return fmt.Errorf(msg)
	}
}

var (
	createPost = pb.AddPostRequest{
		Content: "hello grpc",
		Likes:   500,
	}

	updPost = pb.UpdatePostRequest{
		Id:      1,
		Content: "Hello, gRPC",
		Likes:   100500,
	}

	delPost = pb.DeletePostRequest{Id: 1}

	getPost = pb.GetPostRequest{Id: 1}

	createComment = pb.AddCommentRequest{
		PostId:  1,
		Content: "what is gRPC?",
	}
)

func remoteCallCreatePost(ctx context.Context, client pb.PostServiceClient) error {
	newPost, err := client.CreatePost(ctx, &createPost)
	if err != nil {
		return fmt.Errorf("error after call to CreatePost: %w", err)
	}

	logger.Infof(ctx, "new post: %+v", newPost)
	return nil
}

func remoteCallUpdatePost(ctx context.Context, client pb.PostServiceClient) error {
	updatedPost, err := client.UpdatePost(ctx, &updPost)
	if err != nil {
		return fmt.Errorf("error after call to UpdatePost: %w", err)
	}

	logger.Infof(ctx, "updated post: %+v", updatedPost)
	return nil
}

func remoteCallDeletePost(ctx context.Context, client pb.PostServiceClient) error {
	_, err := client.DeletePost(ctx, &delPost)
	if err != nil {
		return fmt.Errorf("error after call to DeletePost: %w", err)
	}

	logger.Infof(ctx, "post deleted, id: %d", delPost.Id)
	return nil
}

func remoteCallGetPostByID(ctx context.Context, client pb.PostServiceClient) error {
	response, err := client.GetPostByID(ctx, &getPost)
	if err != nil {
		return fmt.Errorf("error after call to DeletePost: %w", err)
	}

	logger.Infof(ctx, "got post: %+v, related comments: %+v", response.Post, response.Comments)
	return nil
}

func remoteCallCreateComment(ctx context.Context, client pb.CommentServiceClient) error {
	newComment, err := client.CreateComment(ctx, &createComment)
	if err != nil {
		return fmt.Errorf("error after call to DeletePost: %w", err)
	}

	logger.Infof(ctx, "new comment: %+v", newComment)
	return nil
}
