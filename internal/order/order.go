package order

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusCreated    = "Created"
	StatusProcessing = "Processing"
	StatusShipped    = "Shipped"
	StatusCanceled   = "Canceled"
	StatusDelivered  = "Delivered"
)

type Order struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ID          uuid.UUID `json:"id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdateddAt  time.Time `json:"updateddAt"`
}

type CrateOrderDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateOrderStatusDto struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
