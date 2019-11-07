
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE articles
  ADD COLUMN body mediumtext NOT NULL,
  ADD COLUMN created datetime,
  ADD COLUMN updated datetime;

UPDATE articles SET created = CURRENT_TIMESTAMP WHERE created IS NULL;
UPDATE articles SET updated = CURRENT_TIMESTAMP WHERE updated IS NULL;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE articles
  DROP COLUMN body,
  DROP COLUMN created,
  DROP COLUMN updated;
