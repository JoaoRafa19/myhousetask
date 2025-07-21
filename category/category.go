package category

import (
	"context"
	"database/sql"
	"fmt"

	pb "JoaoRafa19/myhousetask/proto"
)

type CategoryServiceServer struct {
	pb.UnimplementedCategoryServiceServer
	db *sql.DB
}

func NewCategoryServiceServer(db *sql.DB) *CategoryServiceServer {
	return &CategoryServiceServer{db: db}
}

func (svc *CategoryServiceServer) Save(ctx context.Context, c *pb.Category) (*pb.CategoryResponse, error) {
	_, err := svc.db.ExecContext(ctx, "INSERT INTO categories (id, name, description, is_active) VALUES (?, ?, ?, ?)", c.Id, c.Name, c.Description, c.IsActive)
	if err != nil {
		return nil, fmt.Errorf("error saving category: %w", err)
	}
	return &pb.CategoryResponse{
		Category: c,
	}, nil
}

func (svc *CategoryServiceServer) Find(ctx context.Context, c *pb.CategoryFilterRequest) (*pb.CategoryListResponse, error) {
	rows, err := svc.db.QueryContext(ctx, "SELECT id, name, description, is_active FROM categories WHERE ? = ?", c.Field, c.Value)
	if err != nil {
		return nil, fmt.Errorf("error finding category: %w", err)
	}
	defer rows.Close()

	categories := []*pb.Category{}
	for rows.Next() {
		var id int32
		var name, description string
		var isActive bool
		err := rows.Scan(&id, &name, &description, &isActive)
		if err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}

		categories = append(categories, &pb.Category{
			Id:          id,
			Name:        name,
			Description: description,
			IsActive:    isActive,
		})
	}
	return &pb.CategoryListResponse{Categories: categories}, nil
}
