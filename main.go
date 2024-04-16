package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	// Update
	gw "panda-server/gen/types" // Update

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	gw.UnimplementedHelloServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) Echo(ctx context.Context, in *gw.StringMessage) (*gw.StringMessage, error) {
	return &gw.StringMessage{Value: "hello world"}, nil
}
func (s *server) SayHello(ctx context.Context, in *gw.HelloRequest) (*gw.HelloReply, error) {
	return &gw.HelloReply{Message: "hello this msg"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8091")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	s := grpc.NewServer()
	gw.RegisterHelloServiceServer(s, &server{})
	// gRPC-Gateway mux
	gwmux := runtime.NewServeMux()
	dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = gw.RegisterHelloServiceHandlerFromEndpoint(context.Background(), gwmux, "127.0.0.1:8091", dops)
	if err != nil {
		log.Fatalln("Failed to register gwmux:", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	gwServer := &http.Server{
		Addr:    "127.0.0.1:8091",
		Handler: grpcHandlerFunc(s, mux),
	}
	log.Println("Serving on http://127.0.0.1:8091")
	log.Fatalln(gwServer.Serve(lis))
}
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
