
-- +migrate Up
CREATE TABLE comments (
  comment_id   BIGSERIAL PRIMARY KEY,
  blog_id      BIGINT       NOT NULL,
  client_id    VARCHAR(255)     NULL,
  user_id      BIGINT           NULL,
  content      TEXT         NOT NULL,
  is_edited    BOOLEAN      NOT NULL DEFAULT FALSE,
  is_deleted   BOOLEAN      NOT NULL DEFAULT FALSE,
  thread_id    VARCHAR(255)     NULL,
  created      BIGINT       NOT NULL,
  modified     BIGINT       NOT NULL,
  CONSTRAINT fk_comments_blog
    FOREIGN KEY (blog_id)
    REFERENCES blogs (id)
    ON DELETE CASCADE
);

CREATE INDEX idx_comments_blog_id
  ON comments (blog_id);

CREATE INDEX idx_comments_blog_created
  ON comments (blog_id, created);

-- +migrate Down
DROP TABLE IF EXISTS comments;
