package daemon

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	repo *alarmRepo
}

func (s *server) SetAlarm(ctx context.Context, req *SetAlarmRequest) (*SetAlarmResponse, error) {
	err := s.repo.SetAlarm(req.Alarm)
	if err != nil {
		return nil, err
	}
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
	return &RemoveAlarmResponse{}, nil
}

func startServer(repo *alarmRepo) error {
	soc, err := net.Listen("tcp", PORT)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	s := &server{repo}
	RegisterAlarmServiceServer(grpcServer, s)
	reflection.Register(grpcServer)
	return grpcServer.Serve(soc)
}
