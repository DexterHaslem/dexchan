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


CREATE TABLE IF NOT EXISTS log_action (
  id        SERIAL PRIMARY KEY,
  auser_id  INT REFERENCES auser (id) NOT NULL,
  action    TEXT,
  entity_id INT                       NOT NULL,
  ts        TIMESTAMP
);

CREATE OR REPLACE FUNCTION add_action_log(ip TEXT, action TEXT, entity_id INT)
  RETURNS VOID AS
$func$
INSERT INTO log_action (auser_id, action, entity_id, ts)
  SELECT
    (SELECT get_auser(ip)),
    $2,
    $3,
    now();
$func$
LANGUAGE SQL;
