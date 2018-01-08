-- dexChan copyright Dexter Haslem <dmh@fastmail.com> 2018
-- This file is part of dexChan
--
-- dexChan is free software: you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.
--
-- dexChan is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU General Public License for more details.
--
-- You should have received a copy of the GNU General Public License
-- along with dexChan. If not, see <http://www.gnu.org/licenses/>.


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

