package dbrepo

import (
	"database/sql"
	"github/somyaranjan99/basic-go-project/internal/repository"
	"github/somyaranjan99/basic-go-project/pkg/config"
)

type RepositoryDBHandler struct {
	App *config.AppConfig
	Db  *sql.DB
}

func NewRepositoryDBHandler(app *config.AppConfig, db *sql.DB) repository.DatabaseRepo {
	return &RepositoryDBHandler{
		App: app,
		Db:  db,
	}
}
