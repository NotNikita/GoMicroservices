package main

import (
	"context"
	"log"
	"os"

	"logger/api"
	"logger/service"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	port = ":7070"
	rpcPort = "5001"
	gRpcPort = "50001"
	MONGO_URL = "MONGO_URL"
)

func main(){
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
	
	logger := api.NewLoggerHandler(logService)
	app := fiber.New()
	api.Routes(app)

	app.Get("/ping", logger.HealthCheck)
	app.Post("/log", logger.CreateLog)
	app.Put("/log", logger.UpdateLog)
	app.Get("/logs", logger.GetAllLogs)
	app.Get("/log/:id", logger.GetLogById)
	app.Get("/logs/drop", logger.ClearLogs)

	log.Fatal(app.Listen(port))
}

func connectToMongoDb(uriPath string) (*mongo.Client) {
	client, err := mongo.Connect(options.Client().ApplyURI(uriPath))
	if err != nil {
		log.Fatalf("Failed to create mongo client: %s", err)
		panic(err)
	}
	return client
}