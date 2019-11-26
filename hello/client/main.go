package main
import (
	pb "../../proto/hello" // 引入proto包
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)
const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
	//  是否开启TLS认证
	OpenTLS = true
)



//  自定义认证
type customCredential struct{}
//实现自定义的认证接口
func (c customCredential) GetRequestMetadata(ctx context.Context, uri...string) (map[string]string, error) {
	return map[string]string{
		"appid":"101010",
		"appkey":"i am key",
	},nil
}
// 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool  {
	return OpenTLS
}







func main() {
	// TLS连接 这里的serverName是在生成pem文件时填写的servername
	creds,err := credentials.NewClientTLSFromFile("keys/server.pem","127.0.0.1:50052")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	// 连接
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	//  conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用方法
	req := &pb.HelloRequest{Name: "gRPC"}
	res, err := c.SayHello(context.Background(), req)
	if err != nil {
		grpclog.Fatalln(err)
	}
	fmt.Println(res)
}