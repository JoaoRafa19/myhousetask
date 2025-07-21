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

	userId := sql.NullString{
		String: id,
	}

	return fs.db.ListFamiliesForUser(ctx, userId)
}
