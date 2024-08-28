package main

import (
	"bash_algorithm/grpcTest/pb"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"log"
	// 确保这里是正确的包路径
)

const (
	address = "localhost:50051"
)

func main() {
	// 连接到 gRPC 服务器
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	c := pb.NewPushServiceClient(conn)

	// 创建订阅请求
	clientID := "client1" // 你可以修改为其他客户端 ID
	stream, err := c.Subscribe(context.Background(), &pb.SubscribeRequest{ClientId: clientID})
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	// 接收服务端推送的消息
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("error receiving message: %v", err)
			}
			fmt.Printf("Received message: %s\n", msg.GetMessage())
		}
	}()

	// 让客户端保持运行状态
	select {}
}
