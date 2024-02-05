package postgresql

import (
	"context"
	"errors"
	"fmt"
	"homework-8/internal/pkg/logger"
	"homework-8/internal/utils"
	"time"

	"github.com/opentracing/opentracing-go"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v4"

	"homework-8/internal/pkg/db"
	"homework-8/internal/pkg/repository"
)

// Repo is a struct representing the PostgreSQL repository for handling Post and Comment entities.
type Repo struct {
	db db.DBops
}

// NewRepo creates a new instance of the PostgreSQL repository.
func NewRepo(database db.DBops) *Repo {
	return &Repo{db: database}
}

const (
	msgErrorExecSQL = "error after executing SQL query"
	msgNoRows       = "no rows found with id"
)

// AddPost adds a new Post to the database and returns its ID.
func (r *Repo) AddPost(ctx context.Context, post *repository.Post) (*repository.Post, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresql: AddPost")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "Repo.AddPost")))

	var id int64
	var createdAt time.Time

	err := r.db.ExecQueryRow(ctx, `INSERT INTO posts(content, likes) VALUES($1, $2) RETURNING id, created_at;`, post.Content, post.Likes).Scan(&id, &createdAt)

	if err != nil {
		return nil, utils.ReportError(ctx, msgErrorExecSQL, err)
	}

	post.CreatedAt = createdAt
	post.ID = id

	return post, err
}

// AddComment adds a new Comment to the database and returns its ID.
func (r *Repo) AddComment(ctx context.Context, comment *repository.Comment) (*repository.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresql: AddComment")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "Repo.AddComment")))

	var id int64
	var createdAt time.Time

	err := r.db.ExecQueryRow(ctx, `INSERT INTO comments(post_id, content) VALUES($1, $2) RETURNING id, created_at;`, comment.PostID, comment.Content).Scan(&id, &createdAt)
	if err != nil {
		return nil, utils.ReportError(ctx, msgErrorExecSQL, err)
	}

	comment.CreatedAt = createdAt
	comment.ID = id
	return comment, nil
}

// UpdatePost updates an existing Post in the database.
func (r *Repo) UpdatePost(ctx context.Context, post *repository.Post) (*repository.Post, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresql: UpdatePost")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "Repo.UpdatePost")))

	var id int64
	var createdAt time.Time

	err := r.db.ExecQueryRow(ctx, `UPDATE posts SET content=$1, likes=$2 WHERE id=$3 RETURNING id, created_at;`, post.Content, post.Likes, post.ID).Scan(&id, &createdAt)
	if errors.Is(err, pgx.ErrNoRows) {
		logger.Errorf(ctx, "no rows with id %d found: %w", post.ID, err)
		return nil, repository.ErrObjectNotFound
	}

	if err != nil {
		return nil, utils.ReportError(ctx, msgErrorExecSQL, err)
	}

	post.CreatedAt = createdAt

	return post, nil
}

// GetPostByID retrieves a Post from the database based on its ID.
func (r *Repo) GetPostByID(ctx context.Context, id int64) (*repository.Post, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresql: GetPostByID")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "Repo.GetPostByID")))

	var post repository.Post
	err := r.db.Get(ctx, &post, "SELECT id, content, likes, created_at FROM posts WHERE id=$1;", id)

	if errors.Is(err, pgx.ErrNoRows) {
		logger.Errorf(ctx, msgNoRows, id, err)
		return nil, utils.ReportError(ctx, fmt.Sprintf("%s: %d", msgNoRows, id), err)
	}

	if err != nil {
		return nil, utils.ReportError(ctx, msgErrorExecSQL, err)
	}
	return &post, nil
}

// GetCommentsByPostID retrieves Comments associated with a specific Post from the database.
func (r *Repo) GetCommentsByPostID(ctx context.Context, postID int64) ([]repository.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresql: GetCommentsByPostID")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "Repo.GetCommentsByPostID")))

	var comments []repository.Comment
	err := r.db.Select(ctx, &comments, "SELECT id, content, created_at FROM comments WHERE post_id=$1;", postID)
	if err != nil {
		return nil, utils.ReportError(ctx, msgErrorExecSQL, err)
	}

	return comments, nil
}

// DeletePostByID deletes a Post from the database based on its ID.
func (r *Repo) DeletePostByID(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "postgresql: DeletePostByID")
	defer span.Finish()

	ctx = logger.ToContext(ctx, logger.FromContext(ctx).With(zap.String("method", "Repo.DeletePostByID")))

	res, err := r.db.Exec(ctx, "DELETE FROM posts WHERE id=$1", id)
	if err != nil {
		return utils.ReportError(ctx, msgErrorExecSQL, err)
	}

	if res.RowsAffected() == 0 {
		return utils.ReportError(ctx, fmt.Sprintf("%s: %d", msgNoRows, id), repository.ErrObjectNotFound)
	}
	return nil
}
