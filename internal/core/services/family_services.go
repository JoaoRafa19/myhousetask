package services

import (
	gen "JoaoRafa19/myhousetask/db/gen"
	"context"
	"database/sql"
)

type FamilyServices struct {
	db *gen.Queries
}

func NewFamilyServices(db *gen.Queries) *FamilyServices {
	return &FamilyServices{db: db}
}

func (fs *FamilyServices) GetFamiliesByUserID(ctx context.Context, id string) ([]gen.Family, error) {

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
		return []gen.Family{}, err
	}

	var fams []gen.Family

	for _, fm := range familiesForUserRow {
		fams = append(fams, gen.Family{
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
