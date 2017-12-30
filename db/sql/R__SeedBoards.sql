INSERT INTO board (name, shortname, description, nsfw)
  SELECT
    'random',
    'r',
    'anything goes',
    TRUE
  WHERE NOT exists(SELECT 1
                   FROM board
                   WHERE shortname = 'r');

INSERT INTO board (name, shortname, description, nsfw)
  SELECT
    'test',
    't',
    'test board',
    TRUE
  WHERE NOT exists(SELECT 1
                   FROM board
                   WHERE shortname = 't');

INSERT INTO attachmenttype (name, ext)
  SELECT
    'webm',
    '.webm'
  WHERE NOT exists(SELECT 1
                   FROM attachmenttype
                   WHERE ext = '.webm');

INSERT INTO attachmenttype (name, ext)
  SELECT
    'jpeg',
    '.jpg'
  WHERE NOT exists(SELECT 1
                   FROM attachmenttype
                   WHERE ext = '.jpg');

INSERT INTO attachmenttype (name, ext)
  SELECT
    'png',
    '.png'
  WHERE NOT exists(SELECT 1
                   FROM attachmenttype
                   WHERE ext = '.png');


INSERT INTO board_attachmenttype (board_id, attachmenttype_id, maxsize_bytes)
  SELECT
    (SELECT id
     FROM board
     WHERE shortname = 'r'),
    (SELECT id
     FROM attachmenttype
     WHERE ext = '.webm'),
    (SELECT 1024 * 8)
  WHERE NOT exists(SELECT 1
                   FROM board_attachmenttype bat
                   WHERE bat.attachmenttype_id = (SELECT id
                                                  FROM attachmenttype
                                                  WHERE ext = '.webm'));

INSERT INTO board_attachmenttype (board_id, attachmenttype_id, maxsize_bytes)
  SELECT
    (SELECT id
     FROM board
     WHERE shortname = 'r'),
    (SELECT id
     FROM attachmenttype
     WHERE ext = '.png'),
    (SELECT 1024 * 4)
  WHERE NOT exists(SELECT 1
                   FROM board_attachmenttype bat
                   WHERE bat.attachmenttype_id = (SELECT id
                                                  FROM attachmenttype
                                                  WHERE ext = '.png'));

INSERT INTO board_attachmenttype (board_id, attachmenttype_id, maxsize_bytes)
  SELECT
    (SELECT id
     FROM board
     WHERE shortname = 'r'),
    (SELECT id
     FROM attachmenttype
     WHERE ext = '.jpg'),
    (SELECT 1024 * 4)
  WHERE NOT exists(SELECT 1
                   FROM board_attachmenttype bat
                   WHERE bat.attachmenttype_id = (SELECT id
                                                  FROM attachmenttype
                                                  WHERE ext = '.jpg'));

