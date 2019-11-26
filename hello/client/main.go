package main
import (
	pb "../../proto/hello" // 引入proto包
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"time"
)
const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
	//  是否开启TLS认证
	OpenTLS = true
)



//  自定义认证
type customCredential struct{}
/*下面两个接口是rpc提供的自定义认证的方法
每次rpc调用都会传输认证信息，customCredential 其实是实现了grpc/credential包内的PerRPCCredentials接
每次调用，token信息会通过请求的metadata传输到服务端
*/
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

// interceptor 客户端拦截器
func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	grpclog.Printf("method=%s req=%v rep=%v duration=%s error=%v\n", method, req, reply, time.Since(start), err)
	return err
}



func main() {
	var err error
	var opts []grpc.DialOption
	if OpenTLS {	//  如果开启了TLS认证
		// TLS连接 这里的serverName是在生成pem文件时填写的servername
		creds,err := credentials.NewClientTLSFromFile("keys/server.pem","127.0.0.1:50052")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts,grpc.WithInsecure())
	}

	//  使用自定义认证
	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	// 指定客户端interceptor
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))


	conn,err := grpc.Dial(Address,opts...)
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