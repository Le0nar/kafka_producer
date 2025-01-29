package handler

import (
	"net/http"

	"github.com/Le0nar/kafka_orders/internal/order"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type service interface {
	CreateOrder(dto order.CrateOrderDto) (uuid.UUID, error)
	UpdateOrderStatus(id uuid.UUID, status string) error
}

type Handler struct {
	service service
}

func NewHandler(s service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var input order.CrateOrderDto

	if err := c.BindJSON(&input); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := h.service.CreateOrder(input)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid": uuid,
	})
}

func isValidStatus(status string) bool {
	switch status {
	case order.StatusCreated, order.StatusProcessing, order.StatusShipped, order.StatusCanceled, order.StatusDelivered:
		return true
	default:
		return false
	}
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	stringedId := c.Param("id")

	id, err := uuid.Parse(stringedId)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	status := c.Query("status")

	if !isValidStatus(status) {
		http.Error(c.Writer, "Invalid status", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateOrderStatus(id, status)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, "resource updated successfully")
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/order", h.CreateOrder)
	r.PATCH("/order/:id/status", h.UpdateOrderStatus)

	return r
}
