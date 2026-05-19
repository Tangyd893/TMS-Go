package service

import (
	"context"
	"log/slog"

	"gorm.io/gorm"
)

type ReportService interface {
	GetDashboard(ctx context.Context) (*DashboardData, error)
	GetOrderStats(ctx context.Context) (*OrderStats, error)
	GetTransportEfficiency(ctx context.Context) (*TransportEfficiency, error)
}

type DashboardData struct {
	PendingDispatch  int64 `json:"pendingDispatch"`
	InTransit        int64 `json:"inTransit"`
	PendingException int64 `json:"pendingException"`
	PendingSettlement int64 `json:"pendingSettlement"`
}

type OrderStats struct {
	TotalOrders int64          `json:"totalOrders"`
	ByStatus    []StatusCount  `json:"byStatus"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type TransportEfficiency struct {
	TotalTasks       int64 `json:"totalTasks"`
	CompletedTasks   int64 `json:"completedTasks"`
	AvgTransitHours  float64 `json:"avgTransitHours"`
}

type reportService struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewReportService(db *gorm.DB, log *slog.Logger) ReportService {
	return &reportService{db: db, log: log}
}

func (s *reportService) GetDashboard(ctx context.Context) (*DashboardData, error) {
	data := &DashboardData{}

	if err := s.db.WithContext(ctx).Table("tms_order").
		Where("deleted_at IS NULL AND status = ?", "pending_dispatch").
		Count(&data.PendingDispatch).Error; err != nil {
		data.PendingDispatch = 0
	}

	if err := s.db.WithContext(ctx).Table("tms_transport_task").
		Where("deleted_at IS NULL AND status IN ?", []string{"departed", "in_transit"}).
		Count(&data.InTransit).Error; err != nil {
		data.InTransit = 0
	}

	if err := s.db.WithContext(ctx).Table("tms_exception").
		Where("deleted_at IS NULL AND status = ?", "pending").
		Count(&data.PendingException).Error; err != nil {
		data.PendingException = 0
	}

	if err := s.db.WithContext(ctx).Table("fin_receivable").
		Where("deleted_at IS NULL AND status = ?", "pending").
		Count(&data.PendingSettlement).Error; err != nil {
		data.PendingSettlement = 0
	}

	return data, nil
}

func (s *reportService) GetOrderStats(ctx context.Context) (*OrderStats, error) {
	stats := &OrderStats{}

	if err := s.db.WithContext(ctx).Table("tms_order").
		Where("deleted_at IS NULL").
		Count(&stats.TotalOrders).Error; err != nil {
		return nil, err
	}

	var counts []StatusCount
	if err := s.db.WithContext(ctx).Table("tms_order").
		Select("status, count(*) as count").
		Where("deleted_at IS NULL").
		Group("status").
		Find(&counts).Error; err != nil {
		return nil, err
	}
	stats.ByStatus = counts

	return stats, nil
}

func (s *reportService) GetTransportEfficiency(ctx context.Context) (*TransportEfficiency, error) {
	eff := &TransportEfficiency{}

	if err := s.db.WithContext(ctx).Table("tms_transport_task").
		Where("deleted_at IS NULL").
		Count(&eff.TotalTasks).Error; err != nil {
		return nil, err
	}

	if err := s.db.WithContext(ctx).Table("tms_transport_task").
		Where("deleted_at IS NULL AND status IN ?", []string{"signed", "closed"}).
		Count(&eff.CompletedTasks).Error; err != nil {
		return nil, err
	}

	var avgHours *float64
	if err := s.db.WithContext(ctx).Table("tms_transport_task").
		Select("EXTRACT(EPOCH FROM AVG(sign_time - depart_time)) / 3600").
		Where("deleted_at IS NULL AND sign_time IS NOT NULL AND depart_time IS NOT NULL").
		Scan(&avgHours).Error; err == nil && avgHours != nil {
		eff.AvgTransitHours = *avgHours
	}

	return eff, nil
}
