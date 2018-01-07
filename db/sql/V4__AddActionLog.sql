CREATE TABLE IF NOT EXISTS log_action (
  id       SERIAL PRIMARY KEY,
  auser_id INT REFERENCES auser (id) NOT NULL,
  action   TEXT,
  ts       TIMESTAMP
);

CREATE OR REPLACE FUNCTION add_action_log(ip TEXT, action TEXT)
  RETURNS VOID AS
$func$
INSERT INTO log_action (auser_id, action, ts)
  SELECT
    (SELECT get_auser(ip)),
    $2,
    now();
$func$
LANGUAGE SQL;
