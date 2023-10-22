
-- +migrate Up
CREATE TABLE `blogs` (
  `id`          INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `author_id`   INT NOT NULL,
  `title`       TEXT NOT NULL,
  `content`     TEXT NOT NULL,
  `description` TEXT NOT NULL,
  `thumbnail_image_file_name` TEXT,
  `is_public`   BOOLEAN NOT NULL DEFAULT 1,
  `created`     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `modified`    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE `tags` (
  `id`          INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name`        VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE blogs_tags (
  id          INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  blog_id     INT NOT NULL,
  tag_id      INT NOT NULL,
  FOREIGN KEY(blog_id) REFERENCES blogs(id) ON DELETE CASCADE,
  FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE,
  UNIQUE(blog_id, tag_id)
);

CREATE TABLE users (
  `id`          INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `name`        TEXT NOT NULL,
  `email`       VARCHAR(255) NOT NULL UNIQUE,
  `password`    TEXT NOT NULL,
  `created`     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  `modified`    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);

-- +migrate Down
