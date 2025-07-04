
-- +migrate Up
CREATE TABLE IF NOT EXISTS privacy_policy (
  id                      VARCHAR(36) PRIMARY KEY,
  content                     TEXT             NULL,
  created                 TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  modified                TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER update_privacy_policy_trigger_mod
BEFORE UPDATE ON privacy_policy
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column_timestamp();

-- +migrate Down
DROP TRIGGER IF EXISTS update_privacy_policy_trigger_mod ON privacy_policy;

DROP TABLE IF EXISTS privacy_policy;
