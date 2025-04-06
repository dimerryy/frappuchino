package service

import (
	"strings"

	"hot-coffee/internal/dal"
)

type ReportService struct {
	Repo dal.ReportRepository
}

func NewReportService(repo dal.ReportRepository) *ReportService {
	return &ReportService{Repo: repo}
}

func (s *ReportService) SearchReports(q string, filter string, minPrice, maxPrice float64) (*dal.SearchResult, error) {
	filters := []string{"all"}
	if filter != "" {
		filters = strings.Split(filter, ",")
	}
	return s.Repo.SearchReports(q, filters, minPrice, maxPrice)
}
