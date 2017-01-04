package main

import (
	pb "github.com/goguardian/aws-xray-go/examples/protobuf/demo"
	"github.com/goguardian/aws-xray-go/xray"
	"log"
	"net"

	"golang.org/x/net/context"
)

const listenAddr = ":2000"
const serviceName = "grpc-server-two"

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := xray.NewGRPCServer(serviceName)

	pb.RegisterDemoServer(server, &grpcServer{})

	log.Println("Listening on:", listenAddr)
	log.Fatal(server.Serve(listener))
}

type grpcServer struct{}

func (g *grpcServer) Hi(
	ctx context.Context,
	req *pb.HiRequest,
) (*pb.HiResponse, error) {

	ctx = xray.NewContext(serviceName, ctx)
	defer xray.Close(ctx)

	http, err := xray.GetHTTPClient(ctx)
	if err != nil {
		panic(err)
	}

	http.Get("http://www.example.com")

	return &pb.HiResponse{Message: "hiya"}, nil
}

func (g *grpcServer) Hello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloResponse, error) {

	return &pb.HelloResponse{Message: "hello"}, nil
}
