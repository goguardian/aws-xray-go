package main

import (
	pb "github.com/goguardian/aws-xray-go/examples/protobuf/demo"
	"github.com/goguardian/aws-xray-go/xray"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	ctx := xray.NewContext("grpc-server-one", nil)
	defer xray.Close(ctx)

	conn, err := xray.GetGRPCClientConn("127.0.0.1:2000",
		grpc.WithInsecure(),
		grpc.WithTimeout(10*time.Second))
	if err != nil {
		panic(err)
	}
	client := pb.NewDemoClient(conn)

	resp, hierr := client.Hi(ctx, &pb.HiRequest{Message: "hola"})
	if hierr != nil {
		panic(hierr)
	}
	log.Println(resp.Message)

	helloResp, err := client.Hello(ctx, &pb.HelloRequest{Message: "hello"})
	if err != nil {
		panic(err)
	}
	log.Println(helloResp.Message)
}
