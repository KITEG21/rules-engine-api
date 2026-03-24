-- +migrate Down
DROP INDEX IF EXISTS idx_rules_name;
DROP INDEX IF EXISTS idx_rules_is_active;
DROP TABLE IF EXISTS rules;
