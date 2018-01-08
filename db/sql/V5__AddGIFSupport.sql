UPDATE board b
SET b.allowed_attachment_exts = concat(b.allowed_attachment_exts, ',.gif');