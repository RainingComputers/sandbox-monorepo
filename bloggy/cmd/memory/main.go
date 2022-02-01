package main

import (
	"log"

	"github.com/gin-gonic/gin"

	v1 "examples/bloggy/pkg/routes/v1"
	"examples/bloggy/pkg/storage"
)

func run() error {
	store := storage.CreateMemoryStore()

	router := gin.Default()
	v1.CreateRoutes(store, router)

	router.Run()

	return nil
}

func main() {
	log.Fatalf(run().Error())
}
