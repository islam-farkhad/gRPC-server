package posts

import (
	"homework-8/internal/pkg/repository"
	pb "homework-8/pkg/posts"
)

// Implementation represents the service implementation for posts.
type Implementation struct {
	Repo repository.Repo
	pb.UnimplementedPostServiceServer
}

// New ctor for the posts service implementation.
func New(repo repository.Repo) *Implementation {
	return &Implementation{Repo: repo}
}
