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
        (SELECT 1024 * 5),
        'webm,png,jpg'
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
    (SELECT 1024 * 6),
    'webm,png,jpg'
  WHERE NOT exists(SELECT 1
                   FROM board
                   WHERE shortname = 't');