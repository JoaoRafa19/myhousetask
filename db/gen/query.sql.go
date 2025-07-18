// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const countFamilies = `-- name: CountFamilies :one
SELECT count(*) FROM families
`

func (q *Queries) CountFamilies(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countFamilies)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTasksCompletedToday = `-- name: CountTasksCompletedToday :one
SELECT count(*) FROM tasks WHERE status = 'completed' AND DATE(completed_at) = CURDATE()
`

func (q *Queries) CountTasksCompletedToday(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTasksCompletedToday)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTasksPending = `-- name: CountTasksPending :one
SELECT count(*) FROM tasks WHERE status = 'pending'
`

func (q *Queries) CountTasksPending(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countTasksPending)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countUsers = `-- name: CountUsers :one
SELECT count(*) FROM users
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createFamily = `-- name: CreateFamily :exec
INSERT INTO families (name, description) VALUES (?, ?)
`

type CreateFamilyParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

func (q *Queries) CreateFamily(ctx context.Context, arg CreateFamilyParams) error {
	_, err := q.db.ExecContext(ctx, createFamily, arg.Name, arg.Description)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, name, email, password_hash)
VALUES (?, ?, ?, ?)
`

type CreateUserParams struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
	)
	return err
}

const dashboardPage = `-- name: DashboardPage :many
SELECT 
    f.id as id_familia, 
    f.name as nome_familia, 
    DATE_FORMAT(f.created_at, '%Y-%m-%d %H:%i:%s') as created_at, 
    f.is_active as status,
    COUNT(fm.id) as total_membros
from families f 
LEFT JOIN family_members fm on fm.family_id = f.id 
LEFT JOIN users u on fm.user_id = u.id
GROUP BY f.id, f.name, f.created_at, f.is_active
ORDER BY f.created_at DESC
LIMIT 5
`

type DashboardPageRow struct {
	IDFamilia    int32        `json:"id_familia"`
	NomeFamilia  string       `json:"nome_familia"`
	CreatedAt    string       `json:"created_at"`
	Status       sql.NullBool `json:"status"`
	TotalMembros int64        `json:"total_membros"`
}

func (q *Queries) DashboardPage(ctx context.Context) ([]DashboardPageRow, error) {
	rows, err := q.db.QueryContext(ctx, dashboardPage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DashboardPageRow
	for rows.Next() {
		var i DashboardPageRow
		if err := rows.Scan(
			&i.IDFamilia,
			&i.NomeFamilia,
			&i.CreatedAt,
			&i.Status,
			&i.TotalMembros,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastFiveFamilies = `-- name: GetLastFiveFamilies :many
SELECT id, name, created_at, is_active, description FROM families ORDER BY created_at DESC LIMIT 5
`

func (q *Queries) GetLastFiveFamilies(ctx context.Context) ([]Family, error) {
	rows, err := q.db.QueryContext(ctx, getLastFiveFamilies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Family
	for rows.Next() {
		var i Family
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.IsActive,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, created_at FROM users WHERE email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
	)
	return i, err
}

const getWeeklyTaskCompletionStats = `-- name: GetWeeklyTaskCompletionStats :many
SELECT 
    DATE(completed_at) as completion_date,
    COUNT(*) as completed_count
FROM tasks
WHERE 
    status = 'completed' AND 
    completed_at >= CURDATE() - INTERVAL 7 DAY
GROUP BY completion_date
ORDER BY completion_date DESC
`

type GetWeeklyTaskCompletionStatsRow struct {
	CompletionDate time.Time `json:"completion_date"`
	CompletedCount int64     `json:"completed_count"`
}

func (q *Queries) GetWeeklyTaskCompletionStats(ctx context.Context) ([]GetWeeklyTaskCompletionStatsRow, error) {
	rows, err := q.db.QueryContext(ctx, getWeeklyTaskCompletionStats)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWeeklyTaskCompletionStatsRow
	for rows.Next() {
		var i GetWeeklyTaskCompletionStatsRow
		if err := rows.Scan(&i.CompletionDate, &i.CompletedCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRecentFamilies = `-- name: ListRecentFamilies :many
SELECT 
    f.id as id_familia, 
    f.name as nome_familia, 
    DATE_FORMAT(f.created_at, '%Y-%m-%d %H:%i:%s') as created_at, 
    f.is_active as status,
    COUNT(fm.id) as total_membros
from families f 
LEFT JOIN family_members fm on fm.family_id = f.id 
LEFT JOIN users u on fm.user_id = u.id
GROUP BY f.id, f.name, f.created_at, f.is_active
ORDER BY f.created_at DESC
LIMIT 5
`

type ListRecentFamiliesRow struct {
	IDFamilia    int32        `json:"id_familia"`
	NomeFamilia  string       `json:"nome_familia"`
	CreatedAt    string       `json:"created_at"`
	Status       sql.NullBool `json:"status"`
	TotalMembros int64        `json:"total_membros"`
}

func (q *Queries) ListRecentFamilies(ctx context.Context) ([]ListRecentFamiliesRow, error) {
	rows, err := q.db.QueryContext(ctx, listRecentFamilies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRecentFamiliesRow
	for rows.Next() {
		var i ListRecentFamiliesRow
		if err := rows.Scan(
			&i.IDFamilia,
			&i.NomeFamilia,
			&i.CreatedAt,
			&i.Status,
			&i.TotalMembros,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
