-- +migrate Up
ALTER TABLE summaries
    ADD COLUMN deleted_at TIMESTAMP(6) NULL DEFAULT NULL COMMENT '削除日時（論理削除）' AFTER updated_at,
    ADD INDEX idx_deleted_at (deleted_at);

ALTER TABLE users
    ADD COLUMN deleted_at TIMESTAMP(6) NULL DEFAULT NULL COMMENT '削除日時（論理削除）' AFTER updated_at,
    ADD INDEX idx_deleted_at (deleted_at);

-- +migrate Down
ALTER TABLE summaries
    DROP INDEX idx_deleted_at,
    DROP COLUMN deleted_at;

ALTER TABLE users
    DROP INDEX idx_deleted_at,
    DROP COLUMN deleted_at;
