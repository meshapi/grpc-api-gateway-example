package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	gorillaws "github.com/gorilla/websocket"
	"github.com/meshapi/grpc-api-gateway-examples/gen"
	"github.com/meshapi/grpc-api-gateway/gateway"
	"github.com/meshapi/grpc-api-gateway/websocket"
	"github.com/meshapi/grpc-api-gateway/websocket/wrapper/gorillawrapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":40000")
	if err != nil {
		log.Fatalf("failed to bind: %s", err)
	}

	server := grpc.NewServer()
	gen.RegisterUserServiceServer(server, NewUserService())
	gen.RegisterChatServiceServer(server, NewChatService())
	reflection.Register(server)

	connection, err := grpc.NewClient(":40000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %s", err)
	}

	upgrader := gorillaws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	upgradeFunc := func(w http.ResponseWriter, r *http.Request) (websocket.Connection, error) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("ws error: %s", err)
			return nil, fmt.Errorf("failed to upgrade: %w", err)
		}

		return gorillawrapper.New(c), nil
	}
	restGateway := gateway.NewServeMux(
		gateway.WithWebsocketUpgrader(upgradeFunc),
	)
	gen.RegisterChatServiceHandlerClient(context.Background(), restGateway, gen.NewChatServiceClient(connection))
	gen.RegisterUserServiceHandlerClient(context.Background(), restGateway, gen.NewUserServiceClient(connection))

	httpMux := http.NewServeMux()
	httpMux.Handle("/web/", http.FileServer(http.Dir("templates")))
	httpMux.Handle("/", restGateway)

	go func() {
		log.Printf("starting HTTP on port 4000...")
		if err := http.ListenAndServe(":4000", httpMux); err != nil {
			log.Fatalf("failed to start HTTP Rest Gateway service: %s", err)
		}
	}()

	log.Printf("starting gRPC on port 40000...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to start gRPC server: %s", err)
	}
}
