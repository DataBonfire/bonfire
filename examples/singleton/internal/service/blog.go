package service

import (
	"context"

	pb "github.com/databonfire/bonfire/examples/singleton/api/blog/v1"
)

type BlogService struct {
	pb.UnimplementedBlogServer
}

func NewBlogService() *BlogService {
	return &BlogService{}
}

func (s *BlogService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostReply, error) {
	return &pb.CreatePostReply{}, nil
}
func (s *BlogService) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostReply, error) {
	return &pb.UpdatePostReply{}, nil
}
func (s *BlogService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostReply, error) {
	return &pb.DeletePostReply{}, nil
}
func (s *BlogService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostReply, error) {
	return &pb.GetPostReply{}, nil
}
func (s *BlogService) ListPost(ctx context.Context, req *pb.ListPostRequest) (*pb.ListPostReply, error) {
	return &pb.ListPostReply{}, nil
}
