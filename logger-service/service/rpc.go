package service

import "log"

type RPCServer struct {
	logger *LoggerService
}

type RPCPayload struct {
	Name string
	Data string
}

func NewRPCServer(logger *LoggerService) *RPCServer {
	return &RPCServer{
		logger: logger,
	}
}

func (this *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	err := this.logger.Insert(LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	})
	if err != nil {
		log.Printf("logger-service.RPCServer.LogInfo: failure happened %v", err)
	}

	*resp = "Processed payload via RPC: " + payload.Name

	return nil
}
