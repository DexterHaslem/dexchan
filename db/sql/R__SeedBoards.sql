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

DO $$
DECLARE
  boards TEXT [] := ARRAY ['a', 'b', 'c', 'd'];
  b      TEXT;
BEGIN
  FOREACH b IN ARRAY boards LOOP
    RAISE NOTICE '%', b;

    INSERT INTO board
    (name, shortname, description, nsfw, max_attachment_size, allowed_attachment_exts)
      SELECT
        concat(b, ' board' :: TEXT),
        b,
        'board description here',
        FALSE,
        (SELECT 1024 * 1024 * 5),
        '.webm,.png,.jpg'
      WHERE NOT exists(SELECT 1
                       FROM board
                       WHERE shortname = b);

  END LOOP;
END;
$$;

INSERT INTO board
(name, shortname, description, nsfw, max_attachment_size, allowed_attachment_exts)
  SELECT
    'nsfw',
    'n',
    'anything NSFW goes here',
    TRUE,
    (SELECT 1024 * 1024 * 5),
    '.webm,.png,.jpg'
  WHERE NOT exists(SELECT 1
                   FROM board
                   WHERE shortname = 'n');