package main

import (
	pb "../../proto/hello"
	"google.golang.org/grpc/metadata" // 引入grpc meta包

	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	// rpc服务接口中解析metadata中的信息并验证
	md,ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil,grpc.Errorf(codes.Unauthenticated,"无Token认证信息")
	}

	var (
		appid  string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}
	if appid != "101010" || appkey != "i am key" {
		return nil, grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s. \n Token info: appid=%s,appkey=%s", in.Name, appid, appkey)
	return resp, nil
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
	fmt.Println("Listen on " + Address+" with TLS + TOKEN")
	s.Serve(listen)
}
