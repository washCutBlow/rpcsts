package main

import (
	pb "../../proto/hello"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"net"
)

const (
	// gRpc服务地址
	Address = "127.0.0.1:50052"

)

//  定义helloService 并实现约定的接口
type helloService struct{}

func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error){
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello hahfds %s", in.Name)
	return resp,nil
}


// 
var HelloService  = helloService{}

func main()  {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// TLS认证
	//  注意这里的路径的写法，因为已经在goland中设置了当前工作目录是rpcsts
	creds,err := credentials.NewServerTLSFromFile("keys/server.pem", "keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}

	//  实例化grpc server并开启TLS认证
	s := grpc.NewServer(grpc.Creds(creds))
	// 实例化grpc Server
	// s := grpc.NewServer()


	// 注册HelloService
	pb.RegisterHelloServer(s, HelloService)
	fmt.Println("Listen on " + Address+"with TLS")
	s.Serve(listen)
}
