package utils

import (
	"homework-8/internal/pkg/repository"
	pb "homework-8/pkg/posts"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertRepoPostToProtoPost converts RepoPost to ProtoPost
func ConvertRepoPostToProtoPost(post *repository.Post) *pb.Post {
	return &pb.Post{
		Id:        post.ID,
		Content:   post.Content,
		Likes:     post.Likes,
		CreatedAt: timestamppb.New(post.CreatedAt), // convert time.Time to proto Timestamp
	}
}

// ConvertRepoCommentToProtoComment converts RepoComment to ProtoComment
func ConvertRepoCommentToProtoComment(comment repository.Comment) *pb.Comment {
	return &pb.Comment{
		Id:        comment.ID,
		PostId:    comment.PostID,
		Content:   comment.Content,
		CreatedAt: timestamppb.New(comment.CreatedAt),
	}
}
