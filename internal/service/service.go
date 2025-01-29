package service

import (
	"errors"
	"time"

	"github.com/Le0nar/kafka_orders/internal/order"
	"github.com/google/uuid"
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
	return order.ID, nil
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
	// TODO: get order from Kafka
	order := order.Order{ID: id, Status: order.StatusCreated}

	if !isValidTransition(order.Status, status) {
		return errors.New("invalid change of order status")
	}
	order.Status = status

	// TODO: send  order from Kafka

	return nil
}
