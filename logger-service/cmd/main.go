package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"logger/api"
	"logger/logs"
	"logger/service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
)

const (
	port      = ":7070"
	rpcPort   = "5001"
	gRpcPort  = "50001"
	MONGO_URL = "MONGO_URL"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	dsn := os.Getenv(MONGO_URL)
	log.Println("mongodb dsn is", dsn)
	client := connectToMongoDb(dsn)
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect mongo client: %s", err)
			panic(err)
		}
	}()

	// Init logger service
	logService, err := service.NewLoggerService(client)
	if err != nil {
		log.Fatalf("Failed creating logger service %s", err)
	} else {
		log.Println("Logger service started")
	}
	// Register and Init RPC Server
	rpcServer := service.NewRPCServer(logService)
	err = rpc.Register(rpcServer)
	go rpcListen()

	// Init gRPC Server
	grpcServer := service.NewGRPCLogServer(logService)
	go gRPCListen(grpcServer)

	logger := api.NewLoggerHandler(logService)
	app := fiber.New()
	api.Routes(app, logger)

	log.Fatal(app.Listen(port))
}

func connectToMongoDb(uriPath string) *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI(uriPath))
	if err != nil {
		log.Fatalf("Failed to create mongo client: %s", err)
		panic(err)
	}
	return client
}

func rpcListen() error {
	log.Println("Starting RPC server on port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		log.Println("RPC server accepting connections")
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func gRPCListen(grpcServer *service.GRPCLogServer) error {
	log.Println("Starting gRPC server on port", gRpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", gRpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, grpcServer)

	log.Println("Managed to start gRPC server")

	if err := s.Serve(listen); err != nil {
		log.Fatal("Failed to start gRPC server")
		return err
	}
	return nil
}
