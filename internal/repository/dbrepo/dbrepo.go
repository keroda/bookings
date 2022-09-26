package dbrepo

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keroda/bookings/internal/config"
	"github.com/keroda/bookings/internal/repository"
)

type postgresRepo struct {
	App *config.AppConfig
	DB  *pgxpool.Pool
}

func NewPostgresRepo(conn *pgxpool.Pool, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresRepo{
		App: a,
		DB:  conn,
	}
}
