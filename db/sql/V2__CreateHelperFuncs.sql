-- gets an anon-user id for this ip, creates if doesnt exist
-- if exists, updates last seen
CREATE OR REPLACE FUNCTION get_auser(ip TEXT)
  RETURNS INT AS $$
DECLARE existing_id INT := NULL;
BEGIN

  SELECT id
  FROM auser a
  WHERE a.ip = $1 :: INET
  INTO existing_id;

  IF existing_id IS NOT NULL
  THEN
    UPDATE auser a
    SET last_seen = now()
    WHERE a.id = existing_id;
  ELSE
    WITH newr AS (
      INSERT INTO auser (ip, first_seen, last_seen) VALUES (
        $1 :: INET,
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

CREATE OR REPLACE FUNCTION create_thread(bid INT, sub TEXT, descrp TEXT, ip TEXT, attachmentid INT)
  RETURNS INT AS $$
DECLARE
  thread_id INT := NULL;
  poster_id INT := NULL;
BEGIN
  SELECT get_auser(ip)
  INTO poster_id;

  WITH newr AS (
    INSERT INTO thread (subject, description, created_by_id, created_at, hidden)
    VALUES (sub, descrp, poster_id, now(), FALSE)
    RETURNING id)
  SELECT id
  FROM newr
  INTO thread_id;

  INSERT INTO board_thread (board_id, thread_id) VALUES (
    bid, thread_id
  );

  RETURN thread_id;
END;
$$
LANGUAGE plpgsql;

-- since we have to create attachments separate from posts, this cant be atomic and needs
-- to be controlled in our application layer best as possible
CREATE OR REPLACE FUNCTION create_post(tid INT, contents TEXT, ip TEXT)
  RETURNS UUID AS $$
DECLARE
  post_id   INT := NULL;
  poster_id INT := NULL;
BEGIN
  SELECT get_auser(ip)
  INTO poster_id;

  WITH newr AS (
    INSERT INTO post (content, posted_at, posted_by_id, hidden) VALUES
      (contents, now(), poster_id, FALSE)
    RETURNING id)
  SELECT id
  FROM newr
  INTO post_id;

  INSERT INTO thread_post (thread_id, post_id)
  VALUES (tid, post_id);

  RETURN post_id;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_threads(board_id INT)
  RETURNS SETOF THREAD AS $$
BEGIN
  RETURN QUERY SELECT t.*
               FROM board_thread bt
                 INNER JOIN thread t ON bt.thread_id = t.id
               WHERE bt.board_id = $1;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_threads_with_attachments(board_id INT)
  RETURNS TABLE(
    thread_id           INT,
    thread_subject      TEXT,
    thread_description  TEXT,
    created_at          TIMESTAMP,
    hidden              BOOLEAN,
    attachment_origname TEXT,
    attachment_tn_loc   TEXT,
    attachment_loc      TEXT
  ) AS $$
BEGIN
  RETURN QUERY
  SELECT
    ts.id                         AS thread_id,
    ts.subject                    AS thread_subject,
    ts.description                AS thread_description,
    ts.created_at                 AS created_at,
    ts.hidden                     AS hidden,
    coalesce(a.orig_filename, '') AS attachment_origname,
    coalesce(a.tn_location, '')   AS attachment_tn_loc,
    coalesce(a.location, '')      AS attachment_loc
  FROM get_threads($1) ts
    INNER JOIN thread_attachment ta ON ta.thread_id = ts.id
    LEFT JOIN attachment a ON ta.attachment_id = a.id;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_posts(thread_id INT)
  RETURNS SETOF POST AS $$
BEGIN
  RETURN QUERY SELECT p.*
               FROM thread_post tp
                 INNER JOIN post p ON tp.post_id = p.id
               WHERE tp.thread_id = $1;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_posts_with_attachments(thread_id INT)
  RETURNS TABLE(
    post_id             INT,
    post_content        TEXT,
    posted_at           TIMESTAMP,
    hidden              BOOLEAN,
    attachment_origname TEXT,
    attachment_tn_loc   TEXT,
    attachment_loc      TEXT
  ) AS $$
BEGIN
  RETURN QUERY
  SELECT
    ps.id                         AS post_id,
    ps.content                    AS post_content,
    ps.posted_at                  AS posted_at,
    ps.hidden                     AS hidden,
    coalesce(a.orig_filename, '') AS attachment_origname,
    coalesce(a.tn_location, '')   AS attachment_tn_loc,
    coalesce(a.location, '')      AS attachment_loc
  FROM get_posts($1) ps
    INNER JOIN post_attachment pa ON pa.post_id = ps.id
    LEFT JOIN attachment a ON pa.attachment_id = a.id;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION assign_attachment_post(pid INT, attachment_id INT)
  RETURNS VOID AS $$
BEGIN
  DELETE FROM post_attachment pa
  WHERE pa.attachment_id = attachment_id;

  INSERT INTO post_attachment (post_id, attachment_id) VALUES (pid, attachment_id);
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION assign_attachment_thread(tid INT, attachment_id INT)
  RETURNS VOID AS $$
BEGIN
  DELETE FROM thread_attachment ta
  WHERE ta.attachment_id = attachment_id;
  INSERT INTO thread_attachment (thread_id, attachment_id) VALUES (tid, attachment_id);
END;
$$
LANGUAGE plpgsql;
