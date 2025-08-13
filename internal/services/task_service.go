package services

import (
	"JoaoRafa19/myhousetask/store"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type TasksPageData struct {
	Families       []store.Family
	Tasks          []store.GetTasksByFamilyRow
	TaskStatus     []store.TaskStatus
	SelectedFamily store.Family
}

type TaskService struct {
	db *store.Queries
}

func NewTaskService(db *store.Queries) *TaskService {
	return &TaskService{
		db: db,
	}
}

func (s *TaskService) GetTasksByFamily(ctx context.Context, familyID int32) ([]store.GetTasksByFamilyRow, error) {
	familyId := sql.NullInt32{
		Int32: familyID,
		Valid: true,
	}
	return s.db.GetTasksByFamily(ctx, familyId)
}

func (s *TaskService) GetTaskStatus(ctx context.Context) ([]store.TaskStatus, error) {
	return s.db.GetTaskStatus(ctx)
}

func (s *TaskService) CreateTask(ctx context.Context, arg store.CreateTaskParams) (store.Task, error) {
	arg.ID = uuid.New().String()
	err := s.db.CreateTask(ctx, arg)
	if err != nil {
		return store.Task{}, err
	}

	task, err := s.db.GetTask(ctx, arg.ID)
	if err != nil {
		return store.Task{}, err
	}

	return task, nil
}
