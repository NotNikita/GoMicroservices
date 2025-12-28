package service

import (
	"context"
	"logger/logs"
)

type GRPCLogServer struct {
	logs.UnimplementedLogServiceServer
	ls *LoggerService
}

func NewGRPCLogServer(ls *LoggerService) *GRPCLogServer {
	return &GRPCLogServer{
		ls: ls,
	}
}

func (l *GRPCLogServer) LogInfo(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	name := req.GetName()
	data := req.GetData()

	logEntry1 := LogEntry{
		Name: name,
		Data: data,
	}

	err := l.ls.Insert(logEntry1)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "Logged via GRPC",
	}
	return res, nil
}
