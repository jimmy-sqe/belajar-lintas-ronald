package grpc

import (
	"context"

	healthv1 "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/gen/health/v1"
)

type HealthService struct {
	healthv1.UnimplementedHealthServer
}

func NewHealthService() *HealthService { return &HealthService{} }

func (s *HealthService) Check(_ context.Context, req *healthv1.CheckRequest) (*healthv1.CheckReply, error) {
	return &healthv1.CheckReply{Status: "SERVING"}, nil
}
