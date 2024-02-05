package comments

import (
	"homework-8/internal/pkg/repository"
	pb "homework-8/pkg/posts"
)

// Implementation represents the service implementation for comments.
type Implementation struct {
	Repo repository.Repo
	pb.UnimplementedCommentServiceServer
}

// New ctor for the comments service implementation.
func New(repo repository.Repo) *Implementation {
	return &Implementation{Repo: repo}
}
