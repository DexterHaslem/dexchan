INSERT INTO board
(name, shortname, description, nsfw, max_attachment_size, allowed_attachment_exts)
  SELECT
    'random',
    'r',
    'anything goes',
    TRUE,
    (SELECT 1024 * 5),
    'webm,png,jpg'
  WHERE NOT exists(SELECT 1
                   FROM board
                   WHERE shortname = 'r');


INSERT INTO board
(name, shortname, description, nsfw, max_attachment_size, allowed_attachment_exts)
  SELECT
    'test',
    't',
    'for testing',
    TRUE,
    (SELECT 1024 * 6),
    'webm,png,jpg'
  WHERE NOT exists(SELECT 1
                   FROM board
                   WHERE shortname = 'r');