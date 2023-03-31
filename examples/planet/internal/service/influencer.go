package service

import (
	"context"

	pb "github.com/databonfire/bonfire/examples/planet/api/planet/v1"
)

type InfluencerService struct {
	pb.UnimplementedInfluencerServer
}

func NewInfluencerService() *InfluencerService {
	return &InfluencerService{}
}

func (s *InfluencerService) CreateInfluencer(ctx context.Context, req *pb.CreateInfluencerRequest) (*pb.CreateInfluencerReply, error) {
	return &pb.CreateInfluencerReply{}, nil
}
func (s *InfluencerService) UpdateInfluencer(ctx context.Context, req *pb.UpdateInfluencerRequest) (*pb.UpdateInfluencerReply, error) {
	return &pb.UpdateInfluencerReply{}, nil
}
func (s *InfluencerService) DeleteInfluencer(ctx context.Context, req *pb.DeleteInfluencerRequest) (*pb.DeleteInfluencerReply, error) {
	return &pb.DeleteInfluencerReply{}, nil
}
func (s *InfluencerService) GetInfluencer(ctx context.Context, req *pb.GetInfluencerRequest) (*pb.GetInfluencerReply, error) {
	return &pb.GetInfluencerReply{}, nil
}
func (s *InfluencerService) ListInfluencer(ctx context.Context, req *pb.ListInfluencerRequest) (*pb.ListInfluencerReply, error) {
	return &pb.ListInfluencerReply{}, nil
}
