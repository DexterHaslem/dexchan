CREATE TABLE board (
  id          SERIAL PRIMARY KEY,
  name        TEXT    NOT NULL UNIQUE,
  shortname   TEXT    NOT NULL UNIQUE,
  description TEXT,
  nsfw        BOOLEAN NOT NULL
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
  id            SERIAL PRIMARY KEY,
  subject       TEXT                      NOT NULL,
  description   TEXT,
  created_by_id INT REFERENCES auser (id) NOT NULL,
  created_at    TIMESTAMP                 NOT NULL,
  hidden        BOOLEAN                   NOT NULL
);

CREATE TABLE board_thread (
  board_id  INT REFERENCES board (id)  NOT NULL,
  thread_id INT REFERENCES thread (id) NOT NULL,
  PRIMARY KEY (board_id, thread_id)
);

CREATE TABLE post (
  id           SERIAL PRIMARY KEY,
  content      TEXT                      NOT NULL,
  posted_at    TIMESTAMP                 NOT NULL,
  posted_by_id INT REFERENCES auser (id) NOT NULL,
  hidden       BOOLEAN                   NOT NULL
);

CREATE TABLE thread_post (
  thread_id INT REFERENCES thread (id) NOT NULL,
  post_id   INT REFERENCES post (id)   NOT NULL,
  PRIMARY KEY (thread_id, post_id)
);

CREATE TABLE attachmenttype (
  id   SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  -- moved to join so per-board sizes
  -- maxsize_bytes INT  NOT NULL,
  ext  TEXT NOT NULL UNIQUE
);

CREATE TABLE board_attachmenttype (
  board_id          INT REFERENCES board (id)          NOT NULL,
  attachmenttype_id INT REFERENCES attachmenttype (id) NOT NULL,
  maxsize_bytes     INT                                NOT NULL,
  PRIMARY KEY (board_id, attachmenttype_id)
);

CREATE TABLE attachment (
  id                SERIAL PRIMARY KEY,
  attachmenttype_id INT REFERENCES attachmenttype (id)  NOT NULL,
  orig_filename     TEXT                                NOT NULL,
  uploaded_by_id    INT REFERENCES auser (id)           NOT NULL,
  location          TEXT                                NOT NULL, -- server location for cdn
  tn_location       TEXT                                NOT NULL
);

CREATE TABLE post_attachment (
  post_id       INT REFERENCES post (id)       NOT NULL,
  attachment_id INT REFERENCES attachment (id) NOT NULL,
  PRIMARY KEY (post_id, attachment_id)
);

-- always tn
CREATE TABLE thread_attachment (
  thread_id     INT REFERENCES thread (id)       NOT NULL,
  attachment_id INT REFERENCES attachment (id)   NOT NULL,
  PRIMARY KEY (thread_id, attachment_id)
);

CREATE OR REPLACE VIEW recent_threads_view AS
  SELECT
    b.name    AS board_name,
    t.subject AS subject,
    t.created_at
  FROM post p
    INNER JOIN thread_post tp ON p.id = tp.post_id
    INNER JOIN thread t ON tp.thread_id = t.id
    INNER JOIN board_thread bt ON bt.thread_id = tp.thread_id
    INNER JOIN board b ON bt.board_id = b.id
  ORDER BY t.created_at DESC
  LIMIT 10;

