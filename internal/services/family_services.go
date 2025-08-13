package services

import (
	"JoaoRafa19/myhousetask/store"
	"context"
	"database/sql"
)

type FamilyServices struct {
	db *store.Queries
}

func NewFamilyServices(db *store.Queries) *FamilyServices {
	return &FamilyServices{db: db}
}

func (fs *FamilyServices) GetFamiliesByUserID(ctx context.Context, id string) ([]store.Family, error) {

	user, err := fs.db.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userId := sql.NullString{
		String: user.ID,
		Valid:  true,
	}

	familiesForUserRow, err := fs.db.ListFamiliesForUser(ctx, userId)

	if err != nil {
		return []store.Family{}, err
	}

	var fams []store.Family

	for _, fm := range familiesForUserRow {
		fams = append(fams, store.Family{
			ID:          fm.ID,
			CreatedAt:   fm.CreatedAt,
			Description: fm.Description,
			Name:        fm.Name,
			OwnerID:     fm.OwnerID,
			IsActive:    fm.IsActive,
		})
	}

	return fams, nil
}
