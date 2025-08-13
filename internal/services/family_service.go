package services

import (
	"JoaoRafa19/myhousetask/store"
	"context"
	"database/sql"
)

type FamilyService struct {
	db *store.Queries
}

func NewFamilyService(db *store.Queries) *FamilyService {
	return &FamilyService{
		db: db,
	}
}

func (s *FamilyService) GetFamiliesByUserID(ctx context.Context, userID string) ([]store.Family, error) {
	id := sql.NullString{
		Valid:  true,
		String: userID,
	}

	return s.db.GetFamiliesByUserID(ctx, id)
}
