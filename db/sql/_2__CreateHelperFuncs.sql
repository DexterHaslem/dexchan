CREATE OR REPLACE FUNCTION create_thread(bn TEXT, sub TEXT, descrp TEXT)
  RETURNS INT AS $$
BEGIN

END;
$$ LANGUAGE plpgsql;

-- since we have to create attachments separate from posts, this cant be atomic and needs
-- to be controlled in our application layer best as possible
CREATE OR REPLACE FUNCTION create_post(tid INT, contents TEXT)
  RETURNS UUID AS $$
BEGIN
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION assign_attachments_post(pid INT, attachment_ids INT [])
  RETURNS BOOLEAN AS $$
BEGIN
  RETURN TRUE;
END;
$$ LANGUAGE plpgsql;
