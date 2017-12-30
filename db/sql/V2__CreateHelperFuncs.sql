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
  delete from thread_attachment ta where ta.attachment_id = attachment_id;
  insert into thread_attachment (thread_id, attachment_id) values (tid, attachment_id);
END;
$$
LANGUAGE plpgsql;
