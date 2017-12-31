CREATE TABLE board (
  id                      SERIAL PRIMARY KEY,
  name                    TEXT    NOT NULL UNIQUE,
  shortname               TEXT    NOT NULL UNIQUE,
  description             TEXT,
  nsfw                    BOOLEAN NOT NULL,
  max_attachment_size     INT     NOT NULL,
  allowed_attachment_exts TEXT    NOT NULL
);

CREATE TABLE auser (
  id         SERIAL PRIMARY KEY,
  ip         INET      NOT NULL UNIQUE,
  first_seen TIMESTAMP NOT NULL,
  last_seen  TIMESTAMP NOT NULL
);

CREATE TABLE login (
  id       SERIAL PRIMARY KEY,
  username TEXT    NOT NULL UNIQUE,
  passhash TEXT    NOT NULL,
  enabled  BOOLEAN NOT NULL,
  is_admin BOOLEAN NOT NULL
);

CREATE TABLE login_boardmod (
  user_id  INT REFERENCES login (id)  NOT NULL,
  board_id INT REFERENCES board (id)  NOT NULL,
  enabled  BOOLEAN                    NOT NULL,
  PRIMARY KEY (user_id, board_id)
);

CREATE TABLE thread (
  id                   SERIAL PRIMARY KEY,
  board_id             INT REFERENCES board (id)  NOT NULL,
  subject              TEXT                       NOT NULL,
  description          TEXT,
  created_by_id        INT REFERENCES auser (id)  NOT NULL,
  created_at           TIMESTAMP                  NOT NULL,
  hidden               BOOLEAN                    NOT NULL,
  -- thread and post are two different tables mostly because
  -- thread must always be started with an attachment.
  -- this slightly denormalizes attachment data big woop. needs to be atomic anyway
  attachment_orig_name TEXT                       NOT NULL,
  attachment_tn_loc    TEXT                       NOT NULL,
  attachment_loc       TEXT                       NOT NULL,
  attachment_size      INT                        NOT NULL
);

CREATE TABLE post (
  id                   SERIAL PRIMARY KEY,
  thread_id            INT REFERENCES thread (id) NOT NULL,
  content              TEXT                       NOT NULL,
  posted_at            TIMESTAMP                  NOT NULL,
  posted_by_id         INT REFERENCES auser (id)  NOT NULL,
  hidden               BOOLEAN                    NOT NULL,
  -- had to end up making these not null to not be a pain in the ass to scan
  -- but blank is valid for post. need to add a check to thread to disallow blank
  attachment_orig_name TEXT NOT NULL,
  attachment_tn_loc    TEXT NOT NULL,
  attachment_loc       TEXT NOT NULL,
  attachment_size      INT  NOT NULL
);


CREATE OR REPLACE VIEW recent_threads_view AS
  SELECT
    b.id          AS board_id,
    t.id          AS thread_id,
    t.subject     AS thread_subject,
    t.description AS thread_description,
    t.created_at
  FROM thread t
    INNER JOIN board b ON b.id = t.board_id
  ORDER BY t.created_at DESC
  LIMIT 10;

