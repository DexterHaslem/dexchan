-- gets an anon-user id for this ip, creates if doesnt exist
-- if exists, updates last seen
CREATE OR REPLACE FUNCTION get_auser(ip TEXT)
  RETURNS INT AS $$
DECLARE existing_id INT := NULL;
BEGIN

  SELECT id
  FROM auser a
  WHERE a.ip = $1
  INTO existing_id;

  IF existing_id IS NOT NULL
  THEN
    UPDATE auser a
    SET last_seen = now()
    WHERE a.id = existing_id;
  ELSE
    WITH newr AS (
      INSERT INTO auser (ip, first_seen, last_seen) VALUES (
        $1,
        now(),
        now()
      )
      RETURNING id
    ) SELECT id
      FROM newr
      INTO existing_id;
  END IF;

  RETURN existing_id;
END;
$$
LANGUAGE plpgsql;

