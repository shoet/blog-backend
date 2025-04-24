
-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.modified = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_comments_trigger_mod
BEFORE UPDATE ON comments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column_timestamp();

CREATE OR REPLACE TRIGGER update_user_profile_trigger_mod
BEFORE UPDATE ON user_profile
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column_timestamp();
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER IF EXISTS update_comments_trigger_mod ON comments;
DROP TRIGGER IF EXISTS update_user_profile_trigger_mod ON user_profile;
DROP FUNCTION IF EXISTS update_updated_at_column_timestamp();
