package main

import (
	"github.com/Le0nar/kafka_orders/internal/handler"
	"github.com/Le0nar/kafka_orders/internal/service"
)

func main() {
	service := service.NewService()
	handler := handler.NewHandler(service)

	router := handler.InitRouter()

	router.Run()
}
