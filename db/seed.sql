  WITH data(config, param, payload) AS (
         VALUES ('Develop.mr_robot', 'Database.processing', '{"host": "localhost", "port": "5432", "database": "devdb", "user": "mr_robot", "password": "secret", "schema": "public"}'),
                ('Test.vpn',         'Rabbit.log',          '{"host": "10.0.5.42", "port": "5671", "virtualhost": "/", "user": "guest", "password": "guest"}')
       ),
       inserted_configs AS (
            INSERT INTO configs(name)
            SELECT DISTINCT config
              FROM data
                ON CONFLICT (name) DO NOTHING
         RETURNING name AS config, id AS config_id
       )
INSERT INTO config_params(config_id, name, params)
SELECT c.config_id,
       d.param,
       d.payload::json
  FROM data d
  JOIN inserted_configs c
 USING (config)
