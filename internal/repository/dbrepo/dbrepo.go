package dbrepo

import (
	"database/sql"

	"github.com/msavanan/bookings/internal/config"
	"github.com/msavanan/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

//For testing only
type postgresTestDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(connn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  connn,
	}
}


//For testing only
func NewTestPostgresRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &postgresTestDBRepo{
		App: a,
	}
}