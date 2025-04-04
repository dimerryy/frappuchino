package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hot-coffee/internal/service"
)

type ReportHandler struct {
	Service *service.ReportService
}

func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{Service: service}
}

func (h *ReportHandler) GetSearchReport(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	filter := r.URL.Query().Get("filter")
	minPriceStr := r.URL.Query().Get("minPrice")
	maxPriceStr := r.URL.Query().Get("maxPrice")

	var minPrice, maxPrice float64
	var err error
	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			http.Error(w, "Invalid minPrice", http.StatusBadRequest)
			return
		}
	}
	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			http.Error(w, "Invalid maxPrice", http.StatusBadRequest)
			return
		}
	}

	result, err := h.Service.SearchReports(query, filter, minPrice, maxPrice)
	if err != nil {
		http.Error(w, "Error searching reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
