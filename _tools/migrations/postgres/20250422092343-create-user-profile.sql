
-- +migrate Up
CREATE TABLE user_profile (
  id                BIGSERIAL PRIMARY KEY,
  user_id           BIGINT           NOT NULL UNIQUE, -- PostgreSQLはUNIQUEをつけるとindexは作成される
  profile_image_url TEXT             NULL,
  nickname          TEXT             NULL,
  bio               TEXT             NULL,
  created           TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  modified          TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE user_profile
  ADD CONSTRAINT fk_user_profile_user
  FOREIGN KEY (user_id) REFERENCES users(id);

CREATE OR REPLACE TRIGGER update_user_profile_trigger_mod
BEFORE UPDATE ON user_profile
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_user_profile_trigger_mod ON user_profile;
DROP TABLE IF EXISTS user_profile;

