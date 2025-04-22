
-- +migrate Up

ALTER TABLE comments
DROP COLUMN modified;

ALTER TABLE comments
ADD COLUMN modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;

CREATE OR REPLACE TRIGGER update_comments_trigger_mod
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down

DROP TRIGGER IF EXISTS update_comments_trigger_mod ON comments;

ALTER TABLE comments
DROP COLUMN modified;

ALTER TABLE comments
ADD COLUMN modified BIGINT NOT NULL DEFAULT 0;

