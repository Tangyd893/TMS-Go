package service

import "context"

type Service interface {
	Health(ctx context.Context) Status
	Ready(ctx context.Context) Status
}

type Status struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Health(ctx context.Context) Status {
	return Status{
		Status:  "ok",
		Service: "tms-api",
	}
}

func (s *service) Ready(ctx context.Context) Status {
	return Status{
		Status:  "ready",
		Service: "tms-api",
	}
}
