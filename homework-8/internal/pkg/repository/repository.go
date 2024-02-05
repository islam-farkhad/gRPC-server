//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository

package repository

import "context"

// Repo is an interface representing the repository for posts and comments.
type Repo interface {
	AddPost(ctx context.Context, post *Post) (*Post, error)
	AddComment(ctx context.Context, comment *Comment) (*Comment, error)
	UpdatePost(ctx context.Context, post *Post) (*Post, error)
	GetPostByID(ctx context.Context, id int64) (*Post, error)
	GetCommentsByPostID(ctx context.Context, postID int64) ([]Comment, error)
	DeletePostByID(ctx context.Context, id int64) error
}
