-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY COMMENT 'ユーザーID (UUID)',
    name VARCHAR(255) NOT NULL COMMENT 'ユーザー名',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT 'メールアドレス',
    user_type VARCHAR(50) NOT NULL COMMENT 'ユーザータイプ',
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
    INDEX idx_email (email),
    INDEX idx_user_type (user_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザーテーブル';

CREATE TABLE IF NOT EXISTS summaries (
    id VARCHAR(36) PRIMARY KEY COMMENT 'サマリーID (UUID)',
    title VARCHAR(255) NOT NULL COMMENT 'タイトル',
    description TEXT COMMENT '説明',
    content TEXT NOT NULL COMMENT 'コンテンツ',
    user_id VARCHAR(36) COMMENT 'ユーザーID',
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
    INDEX idx_title (title),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    CONSTRAINT fk_summaries_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='サマリーテーブル';

-- +migrate Down
DROP TABLE IF EXISTS summaries;
DROP TABLE IF EXISTS users;

