package handlers

import (
	db "JoaoRafa19/myhousetask/db/gen"
	"log"
)

type Handler struct {
	db     *db.Queries
	logger *log.Logger
}

func NewHandler(db *db.Queries) *Handler {
	return &Handler{db: db, logger: log.Default()}
}
