package main

import (
	"log"
	"testing"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConfigRepo_GetParams(t *testing.T) {
	db, _ := dbConnect()
	defer db.Close()
	defer db.Exec("TRUNCATE configs CASCADE")

	pgConfigInsert(db, "test1", "param1", `{"x": "foo"}`)
	pgConfigInsert(db, "test1", "param2", `{"y": "foo"}`)

	repo := NewPostgresConfigRepo(db)

	var params = []struct {
		num, configKey, paramKey string
		err                      error
		data                     []byte
	}{
		{"1", "foo", "foo", ErrNotFound, nil},
		{"2", "test1", "param2", nil, []byte(`{"y": "foo"}`)},
		{"3", "notexists", "param2", ErrNotFound, nil},
		{"4", "test1", "notexits", ErrNotFound, nil},
	}

	t.Run("group", func(t *testing.T) {
		for _, d := range params {
			t.Run(d.num, func(t *testing.T) {
				t.Parallel()
				assert := assert.New(t)
				data, err := repo.GetParams(d.configKey, d.paramKey)
				if d.err == nil {
					assert.NoError(err)
					assert.Equal(data, d.data)
				} else {
					assert.Error(err, d.err)
					assert.Nil(data)
				}
			})
		}
	})
}

// pgConfigInsert inserts testing data into the configs/config_params tables.
func pgConfigInsert(db *pgx.ConnPool, config, param, payload string) {
	db.Exec("INSERT INTO configs(name) VALUES ($1)", config)
	_, err := db.Exec(`
		   WITH conf AS (SELECT id FROM configs WHERE name = $1),
		        params(name, payload) AS (VALUES ($2, $3))
		 INSERT INTO config_params(config_id, name, params)
		 SELECT c.id,
		        cp.name,
		        cp.payload::json
		   FROM conf c, params cp
	`, config, param, payload)
	if err != nil {
		log.Panic(err)
	}
}
