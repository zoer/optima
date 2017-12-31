package main

import (
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

var _ ConfigRepo = (*PostgresConfigRepo)(nil)

type PostgresConfigRepo struct {
	db *pgx.ConnPool
}

// NewPostgresConfigRepo allocates PostgreSQL config reposistory.
func NewPostgresConfigRepo(db *pgx.ConnPool) *PostgresConfigRepo {
	return &PostgresConfigRepo{db: db}
}

// GetParams loads config params with given keys from DB.
func (r *PostgresConfigRepo) GetParams(config, param string) ([]byte, error) {
	var data pgtype.JSON

	err := r.db.QueryRow(`
		SELECT cp.params
		  FROM configs c
		  JOIN config_params cp
		    ON c.id = cp.config_id
		 WHERE c.name = $1
		   AND cp.name = $2
		 LIMIT 1
	`, config, param).Scan(&data)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return data.Bytes, nil
}
