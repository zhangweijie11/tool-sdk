package proto

import (
	"context"
	"errors"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	"net"

	pb "gitlab.example.com/zhangweijie/tool-sdk/services/callback/proto"

	"google.golang.org/grpc"
)

type CallbackService struct {
	pb.UnimplementedCallbackServiceServer
}

func (cs *CallbackService) PushProgress(ctx context.Context, in *pb.PushProgressRequest) (*pb.PushProgressResponse, error) {
	return &pb.PushProgressResponse{Code: 200, Msg: "success"}, nil
}

func (cs *CallbackService) PushResult(ctx context.Context, in *pb.PushResultRequest) (*pb.PushResultResponse, error) {
	return &pb.PushResultResponse{Code: 200, Msg: "success"}, nil
}

func InitgRPC() error {
	listen, err := net.Listen("http", global.Config.Server.RPCAddress)
	if err != nil {
		return errors.New(schemas.RPCConnectErr)
	}
	s := grpc.NewServer()
	pb.RegisterCallbackServiceServer(s, &CallbackService{})
	if err = s.Serve(listen); err != nil {
		return errors.New(schemas.RPCConnectErr)
	}
	return err
}
