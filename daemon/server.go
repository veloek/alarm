package daemon

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	repo    *alarmRepo
	updates chan<- []*Alarm
}

func newServer(repo *alarmRepo, updates chan<- []*Alarm) *server {
	s := &server{repo, updates}
	return s
}

func (s *server) SetAlarm(ctx context.Context, req *SetAlarmRequest) (*SetAlarmResponse, error) {
	err := s.repo.SetAlarm(req.Alarm)
	if err != nil {
		return nil, err
	}
	s.notify()
	return &SetAlarmResponse{}, nil
}

func (s *server) GetAlarms(ctx context.Context, req *GetAlarmsRequest) (*GetAlarmsResponse, error) {
	alarms, err := s.repo.GetAlarms()
	if err != nil {
		return nil, err
	}
	return &GetAlarmsResponse{Alarms: alarms}, nil
}

func (s *server) RemoveAlarm(ctx context.Context, req *RemoveAlarmRequest) (*RemoveAlarmResponse, error) {
	err := s.repo.RemoveAlarm(int(req.AlarmId))
	if err != nil {
		return nil, err
	}
	s.notify()
	return &RemoveAlarmResponse{}, nil
}

func (s *server) notify() {
	a, err := s.repo.GetAlarms()
	if err == nil {
		s.updates <- a
	}
}

func (s *server) start() {
	soc, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on port %s: %v\n", PORT, err)
		return
	}
	s.notify() // Initial update.

	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()
	RegisterAlarmServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	err = grpcServer.Serve(soc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gRPC error: %v\n", err)
	}
}
