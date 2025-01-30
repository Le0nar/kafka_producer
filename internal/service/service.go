package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Le0nar/kafka_producer/internal/order"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateOrder(dto order.CrateOrderDto) (uuid.UUID, error) {
	order := order.Order{
		Name:        dto.Name,
		Description: dto.Description,
		ID:          uuid.New(),
		Status:      order.StatusCreated,
		CreatedAt:   time.Now(),
		UpdateddAt:  time.Now(),
	}

	// TODO: send data to Kafka
	return order.ID, sendToKafka(order)
}

// Проверка, можно ли изменить статус
func isValidTransition(currentStatus, newStatus string) bool {
	// допустимые переходы между статусами
	var validTransitions = map[string][]string{
		order.StatusCreated:    {order.StatusProcessing, order.StatusCanceled},
		order.StatusProcessing: {order.StatusShipped, order.StatusCanceled},
		order.StatusShipped:    {order.StatusDelivered, order.StatusCanceled},
		order.StatusCanceled:   {},
		order.StatusDelivered:  {},
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}
	for _, status := range allowedStatuses {
		if status == newStatus {
			return true
		}
	}
	return false
}

func (s *Service) UpdateOrderStatus(id uuid.UUID, status string) error {
	// TODO: get order from Storage
	order := order.Order{ID: id, Status: order.StatusCreated}

	if !isValidTransition(order.Status, status) {
		return errors.New("invalid change of order status")
	}
	order.Status = status

	return sendToKafka(order)
}

func sendToKafka(order order.Order) error {
	// Создаем подключение к Kafka
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
	})

	// Сериализация заказа в JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// Отправляем сообщение в Kafka
	err = writer.WriteMessages(nil, kafka.Message{
		Value: orderJSON,
	})
	if err != nil {
		return err
	}

	log.Printf("Order sent to Kafka: %v", order)

	return nil
}
