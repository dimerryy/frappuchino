package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type OrderHandler interface {
	PostOrder(w http.ResponseWriter, r *http.Request)
	PutOrderByID(w http.ResponseWriter, r *http.Request)
	DeleteOrderByID(w http.ResponseWriter, r *http.Request)
	GetOrderByID(w http.ResponseWriter, r *http.Request)
	GetAllOrders(w http.ResponseWriter, r *http.Request)
	PostCloseOrder(w http.ResponseWriter, r *http.Request)
	GetNumberOfOrderedItems(w http.ResponseWriter, r *http.Request)
	GetOrderedItemsByPeriod(w http.ResponseWriter, r *http.Request)
	BatchProcessOrders(w http.ResponseWriter, r *http.Request)
}

type orderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *orderHandler {
	return &orderHandler{orderService: orderService}
}

func (h *orderHandler) PostOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed to decode", err.Error(), "no order posted")
		return
	}
	err = h.orderService.PostOrUpdate(newOrder, 0)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		slog.Error("Failed to decode", err.Error(), "no order posted")
		return
	}
	slog.Info("order posted")
	w.WriteHeader(http.StatusCreated)
}

func (h *orderHandler) PutOrderByID(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != 3 {
		RespondWithJson(w, ErrorResponse{Message: "Invalid path"}, http.StatusBadRequest)
		slog.Error("Failed", "wrong params", "no order posted")
		return
	}
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	id, err := strconv.Atoi(pathParam[2])
	if err != nil {
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	orderItem, err := h.orderService.GetOrderItemById(id)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	if orderItem.Status == "active" {
		err = h.orderService.PostOrUpdate(newOrder, id)
		if err != nil {
			if err.Error() == "not found" {
				RespondWithJson(w, ErrorResponse{Message: "Order not found"}, http.StatusNotFound)
				slog.Error("Failed", err.Error(), "no order posted")
				return
			}
			if err.Error() == "bad request" {
				RespondWithJson(w, ErrorResponse{Message: "Bad request"}, http.StatusBadRequest)
				slog.Error("Failed", err.Error(), "no order posted")
				return
			}
			RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
			slog.Error("Failed to update", err.Error(), "no order posted")
			return
		}
	} else {
		RespondWithJson(w, ErrorResponse{Message: "order closed"}, http.StatusNotFound)
		slog.Error("Failed", "no order", "order closed")
		return
	}
	slog.Info("order posted", "orderID", id)
}

func (h *orderHandler) DeleteOrderByID(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != 3 {
		RespondWithJson(w, ErrorResponse{Message: "Invalid path"}, http.StatusBadRequest)
		slog.Error("Failed", "wrong params", "no order posted")
		return
	}
	id, err := strconv.Atoi(pathParam[2])
	if err != nil {
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	err = h.orderService.DeleteOrder(id)
	if err != nil {
		if err.Error() == "not found" {
			RespondWithJson(w, ErrorResponse{Message: "Order not found"}, http.StatusNotFound)
			slog.Error("Failed", err.Error(), "no order posted")
			return
		}
		if err.Error() == "bad request" {
			RespondWithJson(w, ErrorResponse{Message: "Bad request"}, http.StatusBadRequest)
			slog.Error("Failed", err.Error(), "no order posted")
			return
		}
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
	}
	slog.Info("order posted", "orderID", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *orderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/orders/"):])
	if err != nil {
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	orderItem, err := h.orderService.GetOrderItemById(id)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	err = setBodyToJson(w, orderItem)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	slog.Info("order posted", "orderID", id)
}

func (h *orderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orderItems, err := h.orderService.GetOrderItem()
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}

	jsonData, err := json.MarshalIndent(orderItems, "", "  ")
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	slog.Info("order posted", "orderID", orderItems)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *orderHandler) PostCloseOrder(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	if len(pathParam) != 4 {
		RespondWithJson(w, ErrorResponse{Message: "Invalid path"}, http.StatusBadRequest)
		slog.Error("Failed", "wrong params", "no order posted")
		return
	}
	id, err := strconv.Atoi(pathParam[2])
	if err != nil {
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	if err := h.orderService.UpdateOrderStatus(id); err != nil {
		if err.Error() == "order is already closed" {
			RespondWithJson(w, ErrorResponse{Message: "Order is already closed"}, http.StatusNotFound)
			slog.Error("Failed", err.Error(), "order is already closed")
			return
		}
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusNotFound)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
	slog.Info("order posted", "orderID", id)
}

func (h *orderHandler) GetNumberOfOrderedItems(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	items, err := h.orderService.GetNumberOfOrderedItems(startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = setBodyToJson(w, items)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no report posted")
		return
	}
}

func (h *orderHandler) GetOrderedItemsByPeriod(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	var result interface{}
	var err error

	switch period {
	case "day":
		result, err = h.orderService.GetOrdersGroupedByDay(month)
	case "month":
		result, err = h.orderService.GetOrdersGroupedByMonth(year)
	default:
		http.Error(w, "Invalid period value", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = setBodyToJson(w, result)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no report posted")
		return
	}
}

func (h *orderHandler) BatchProcessOrders(w http.ResponseWriter, r *http.Request) {
	var req models.BatchOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	resp, err := h.orderService.ProcessBatchOrders(req.Orders)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}

	err = setBodyToJson(w, resp)
	if err != nil {
		RespondWithJson(w, ErrorResponse{Message: err.Error()}, http.StatusInternalServerError)
		slog.Error("Failed", err.Error(), "no order posted")
		return
	}
}
