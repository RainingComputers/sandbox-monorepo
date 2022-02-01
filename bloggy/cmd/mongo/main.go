package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	v1 "examples/bloggy/pkg/routes/v1"
	"examples/bloggy/pkg/storage"
)

func run() error {
	// Create Mongo store
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := storage.CreateMongoStore(ctx, "test", "test")
	if err != nil {
		return err
	}

	// Create routes and start serving
	router := gin.Default()
	v1.CreateRoutes(store, router)

	router.Run()

	return nil
}

func main() {
	log.Fatalf(run().Error())
}
