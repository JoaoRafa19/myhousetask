package services

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"context"
	"database/sql"
	"fmt"
)

type DashboardData struct {
	UserName                 string
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

func (s *DashboardService) GetDashboardData(userId string) (*DashboardData, error) {

	if userId == "" {
		return nil, fmt.Errorf("invalid user id ")
	}

	user, err := s.db.GetUserByID(context.Background(), userId)
	if err != nil {
		return nil, err
	}

	userid := sql.NullString{
		String: userId,
		Valid:  true,
	}

	totalFamilies, err := s.db.CountFamilies(context.Background(), userid)
	if err != nil {
		return nil, err
	}

	totalUsers, err := s.db.CountUsersFamilyMembers(context.Background(), userid)
	if err != nil {
		return nil, err
	}

	totalTasksCompletedToday, err := s.db.CountTasksCompletedToday(context.Background(), userid)
	if err != nil {
		return nil, err
	}

	totalTasksPending, err := s.db.CountTasksPending(context.Background())
	if err != nil {
		return nil, err
	}

	recentFamilies, err := s.db.ListRecentFamilies(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println(totalFamilies, totalUsers, totalTasksCompletedToday, totalTasksPending, recentFamilies)

	return &DashboardData{
		UserName:                 user.Name,
		TotalFamilies:            totalFamilies,
		TotalUsers:               totalUsers,
		TotalTasksCompletedToday: totalTasksCompletedToday,
		TotalTasksPending:        totalTasksPending,
		RecentFamilies:           recentFamilies,
	}, nil
}
