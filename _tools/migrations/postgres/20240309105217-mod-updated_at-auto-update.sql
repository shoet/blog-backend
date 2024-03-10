
-- +migrate Up

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified = EXTRACT(EPOCH FROM CURRENT_TIMESTAMP);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_blogs_trigger_mod
BEFORE UPDATE ON blogs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_trigger_mod
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER IF EXISTS update_blogs_trigger_mod ON blogs;
DROP TRIGGER IF EXISTS update_users_trigger_mod ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
