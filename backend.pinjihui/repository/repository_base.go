package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/op/go-logging"
)

type BaseRepository struct {
	db  *sqlx.DB
	log *logging.Logger
}
