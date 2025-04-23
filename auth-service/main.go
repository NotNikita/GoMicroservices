package main

import (
	"auth/api"
	"auth/service"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const (
	webPort = ":9090"
	DB_URL  = "POSTGRES_DB_URL"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	dsn := os.Getenv(DB_URL)
	log.Println("dsn is", dsn)
	pgPool, err := createDBPool(dsn)
	if err != nil {
		log.Fatalf("Failed connecting to DB %s", err)

	}

	// Init auth service
	userService, err := service.NewUserService(pgPool)
	if err != nil {
		log.Fatalf("Failed creating auth user service %s", err)
	}

	// init controller
	userController := api.NewUserHandler(userService)

	log.Println("Auth service started")
	app := fiber.New()
	api.ApplyMiddleware(app)

	app.Get("/ping", userController.Health)

	app.Post("/login", userController.Login)

	log.Fatal(app.Listen(webPort))
}

/*
* pgx.Connect - Creates a single connection to the db.
*	Suitable for simple app where only 1 connection is needed.
*	Not ideal for production applications where multiple concurrent connections are required.

* pgx.New - Creates a connection pool that manages multiple connections to the database.
* 	The pool automatically manages connections (e.g., reuses idle connections, creates new ones as needed).
* 	More efficient for high-concurrency applications.
 */
func createDBPool(dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database connection string is empty")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Printf("Couldn't open connection Pool: %s", err)
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("Failed to ping database: %s", err)
		return nil, err
	}

	log.Println("Successfully connected to database")
	return pool, nil
}
