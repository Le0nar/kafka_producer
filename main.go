package main

import (
	"github.com/Le0nar/kafka_producer/internal/handler"
	"github.com/Le0nar/kafka_producer/internal/service"
)

func main() {
	service := service.NewService()
	handler := handler.NewHandler(service)

	router := handler.InitRouter()

	router.Run()
}
