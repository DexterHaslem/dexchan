UPDATE board b
SET allowed_attachment_exts = concat(b.allowed_attachment_exts, ',.gif');