package util

import (
	"context"
	"encoding/json"
	"errors"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/middleware/schemas"
	pb "gitlab.example.com/zhangweijie/tool-sdk/services/callback/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

// CallbackgRPC 通过gRPC 形式回调
func CallbackgRPC(validParams interface{}) error {
	var (
		validUrl     string
		validMessage interface{}
	)

	switch validParams.(type) {
	case *global.Progress:
		validUrl = validParams.(*global.Progress).ProgressUrl
		validMessage = &pb.PushProgressRequest{WorkUUID: validParams.(*global.Progress).WorkUUID, ServerName: global.Config.Server.ServerName, Progress: validParams.(*global.Progress).Progress}
	case *global.Result:
		validUrl = validParams.(*global.Result).CallbackUrl
		// 将JSON对象编码为JSON字符串
		jsonData, err := json.Marshal(validParams.(*global.Result).Result)
		if err != nil {
			return errors.New(schemas.JsonParseErr)
		}
		validMessage = &pb.PushResultRequest{WorkUUID: validParams.(*global.Result).WorkUUID, ServerName: global.Config.Server.ServerName, Result: string(jsonData)}
	default:
		return errors.New(schemas.WorkCallbackErr)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(validUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return errors.New(schemas.RPCConnectErr)
	}

	client := pb.NewCallbackServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch validMessage.(type) {
	case *pb.PushProgressRequest:
		response, err := client.PushProgress(ctx, validMessage.(*pb.PushProgressRequest))
		if err != nil {
			return errors.New(schemas.RPCPushErr)
		}

		if response.GetCode() != 200 {
			return errors.New(schemas.RPCPushErr)
		}
	case *pb.PushResultRequest:
		response, err := client.PushResult(ctx, validMessage.(*pb.PushResultRequest))
		if err != nil {
			return errors.New(schemas.RPCPushErr)
		}

		if response.GetCode() != 200 {
			return errors.New(schemas.RPCPushErr)
		}
	}

	return err
}
