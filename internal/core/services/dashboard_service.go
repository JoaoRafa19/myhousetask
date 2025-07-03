package services

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"context"
	"fmt"
)

type DashboardData struct {
	TotalFamilies            int64
	TotalUsers               int64
	TotalTasksCompletedToday int64
	TotalTasksPending        int64
	RecentFamilies           []db.ListRecentFamiliesRow
}

type DashboardService struct {
	db *db.Queries
}

func NewDashboardService(db *db.Queries) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) GetWeeklyActivity() ([]db.GetWeeklyTaskCompletionStatsRow, error) {
	wekelyActivity, err := s.db.GetWeeklyTaskCompletionStats(context.Background())
	if err != nil {
		return nil, err
	}
	return wekelyActivity, nil
}

func (s *DashboardService) GetDashboardData() (DashboardData, error) {
	totalFamilies, err := s.db.CountFamilies(context.Background())
	if err != nil {
		return DashboardData{}, err
	}

	totalUsers, err := s.db.CountUsers(context.Background())
	if err != nil {
		return DashboardData{}, err
	}

	totalTasksCompletedToday, err := s.db.CountTasksCompletedToday(context.Background())
	if err != nil {
		return DashboardData{}, err
	}

	totalTasksPending, err := s.db.CountTasksPending(context.Background())
	if err != nil {
		return DashboardData{}, err
	}

	recentFamilies, err := s.db.ListRecentFamilies(context.Background())
	if err != nil {
		return DashboardData{}, err
	}

	fmt.Println(totalFamilies, totalUsers, totalTasksCompletedToday, totalTasksPending, recentFamilies)

	return DashboardData{
		TotalFamilies:            totalFamilies,
		TotalUsers:               totalUsers,
		TotalTasksCompletedToday: totalTasksCompletedToday,
		TotalTasksPending:        totalTasksPending,
		RecentFamilies:           recentFamilies,
	}, nil
}
