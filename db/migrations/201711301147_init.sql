-- +goose Up
CREATE TABLE configs(
  id serial PRIMARY KEY,
  name text UNIQUE NOT NULL
);

CREATE TABLE config_params(
  id serial PRIMARY KEY,
  config_id integer NOT NULL REFERENCES configs ON DELETE CASCADE,
  name text NOT NULL,
  params json
);

CREATE UNIQUE INDEX config_params_name_idx
    ON config_params(config_id, name);

-- +goose Down
DROP TABLE config_params, configs;
