
-- +migrate Up
CREATE TABLE `blogs` (
  `id`          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  `author_id`   INTEGER NOT NULL,
  `title`       TEXT NOT NULL,
  `content`     TEXT NOT NULL,
  `description`     TEXT NOT NULL,
  `thumbnail_image_file_name` TEXT,
  `is_public`   BOOLEAN NOT NULL DEFAULT 1,
  `created`     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `modified`    DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE `tags` (
  `id`          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  `name`        TEXT NOT NULL UNIQUE
);

CREATE TABLE blogs_tags (
  id          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  blog_id     INTEGER NOT NULL,
  tag_id      INTEGER NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id) ON DELETE CASCADE,
  FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE,
  UNIQUE(blog_id, tag_id)
);

CREATE TABLE users (
  `id`          INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  `name`        TEXT NOT NULL,
  `email`       TEXT NOT NULL UNIQUE,
  `password`    TEXT NOT NULL,
  `created`     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `modified`    DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- +migrate Down
