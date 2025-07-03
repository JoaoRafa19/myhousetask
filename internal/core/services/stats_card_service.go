package services

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"context"
)

type StatsCardService struct {
	db *db.Queries
}

func NewStatsCardService(db *db.Queries) *StatsCardService {
	return &StatsCardService{db: db}
}

type StatsCardInfo struct {
	TotalFamilies            int64
	TotalUsers               int64
	TotalTasksCompletedToday int64
	TotalTasksPending        int64
}

func (s *StatsCardService) GetStatsCardData(ctx context.Context) (StatsCardInfo, error) {

	totalFamilies, err := s.db.CountFamilies(ctx)

	if err != nil {
		return StatsCardInfo{}, err
	}
	totalUsers, err := s.db.CountUsers(ctx)
	if err != nil {
		return StatsCardInfo{}, err
	}
	totalTasksCompletedToday, err := s.db.CountTasksCompletedToday(ctx)
	if err != nil {
		return StatsCardInfo{}, err
	}
	totalTasksPending, err := s.db.CountTasksPending(ctx)
	if err != nil {
		return StatsCardInfo{}, err
	}

	return StatsCardInfo{
		TotalFamilies:            totalFamilies,
		TotalUsers:               totalUsers,
		TotalTasksCompletedToday: totalTasksCompletedToday,
		TotalTasksPending:        totalTasksPending,
	}, nil
}


