package services

import (
	"JoaoRafa19/myhousetask/store"
	"context"
	"database/sql"
	"fmt"
)

type StatsCardService struct {
	db *store.Queries
}

func NewStatsCardService(db *store.Queries) *StatsCardService {
	return &StatsCardService{db: db}
}

type StatsCardInfo struct {
	TotalFamilies            int64
	TotalMembers             int64
	TotalTasksCompletedToday int64
	TotalTasksPending        int64
}

func (s *StatsCardService) GetStatsCardData(ctx context.Context, user_id string) (*StatsCardInfo, error) {

	if user_id == "" {
		return nil, fmt.Errorf("invalid user id ")
	}

	userid := sql.NullString{
		String: user_id,
		Valid:  true,
	}

	totalFamilies, err := s.db.CountFamilies(ctx, userid)

	if err != nil {
		return nil, err
	}
	totalUsers, err := s.db.CountUsersFamilyMembers(ctx, userid)
	if err != nil {
		return nil, err
	}
	totalTasksCompletedToday, err := s.db.CountTasksCompletedToday(ctx, userid)
	if err != nil {
		return nil, err
	}
	totalTasksPending, err := s.db.CountTasksPending(ctx)
	if err != nil {
		return nil, err
	}

	return &StatsCardInfo{
		TotalFamilies:            totalFamilies,
		TotalMembers:             totalUsers,
		TotalTasksCompletedToday: totalTasksCompletedToday,
		TotalTasksPending:        totalTasksPending,
	}, nil
}
