package main

import (
	"log"
	"net"
	"sync"
	"time"

	pb "bash_algorithm/grpcTest/pb" // 替换为你生成的 pb 包的实际路径
	"google.golang.org/grpc"
)

type pushServer struct {
	pb.UnimplementedPushServiceServer
	clients map[string]chan *pb.PushMessage
	mu      sync.Mutex
}

func NewPushServer() *pushServer {
	return &pushServer{
		clients: make(map[string]chan *pb.PushMessage),
	}
}

// Subscribe 实现客户端订阅逻辑
func (s *pushServer) Subscribe(req *pb.SubscribeRequest, stream pb.PushService_SubscribeServer) error {
	clientID := req.ClientId
	msgChan := make(chan *pb.PushMessage, 10)

	s.mu.Lock()
	s.clients[clientID] = msgChan
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, clientID)
		s.mu.Unlock()
		close(msgChan)
	}()

	// 持续向客户端发送消息
	for msg := range msgChan {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

// pushMessageToAllClients 模拟向所有订阅的客户端发送消息
func (s *pushServer) pushMessageToAllClients(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, msgChan := range s.clients {
		msgChan <- &pb.PushMessage{Message: message}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pushServer := NewPushServer()
	pb.RegisterPushServiceServer(grpcServer, pushServer)

	// 启动一个 goroutine，定时推送消息给所有客户端
	go func() {
		for {
			time.Sleep(5 * time.Second)
			pushServer.pushMessageToAllClients("Hello from server!")
		}
	}()

	log.Println("Server is running on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
