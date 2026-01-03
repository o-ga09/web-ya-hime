-- +migrate Up
-- カテゴリテーブルを作成
CREATE TABLE IF NOT EXISTS categories (
    id CHAR(36) PRIMARY KEY COMMENT 'カテゴリID（UUID）',
    name VARCHAR(100) NOT NULL UNIQUE COMMENT 'カテゴリ名',
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='カテゴリマスタ';

-- サブカテゴリテーブルを作成
CREATE TABLE IF NOT EXISTS subcategories (
    id CHAR(36) PRIMARY KEY COMMENT 'サブカテゴリID（UUID）',
    category_id CHAR(36) NOT NULL COMMENT 'カテゴリID',
    name VARCHAR(100) NOT NULL COMMENT 'サブカテゴリ名',
    created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '作成日時',
    updated_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新日時',
    UNIQUE KEY uk_category_name (category_id, name),
    INDEX idx_category_id (category_id),
    INDEX idx_name (name),
    CONSTRAINT fk_subcategory_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='サブカテゴリマスタ';

-- summariesテーブルにカテゴリとサブカテゴリのIDカラムを追加
ALTER TABLE summaries
    ADD COLUMN category_id CHAR(36) NULL COMMENT 'カテゴリID' AFTER content,
    ADD COLUMN subcategory_id CHAR(36) NULL COMMENT 'サブカテゴリID' AFTER category_id,
    ADD INDEX idx_category_id (category_id),
    ADD INDEX idx_subcategory_id (subcategory_id),
    ADD CONSTRAINT fk_summary_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    ADD CONSTRAINT fk_summary_subcategory FOREIGN KEY (subcategory_id) REFERENCES subcategories(id) ON DELETE SET NULL;

-- 既存のcategoryカラムからデータを移行するためのロジック（必要に応じて手動で実行）
-- 既存のcategoryカラムは削除せず、後方互換性のため残す
-- 将来的に削除する場合は別のマイグレーションで実施

-- +migrate Down
ALTER TABLE summaries
    DROP FOREIGN KEY fk_summary_subcategory,
    DROP FOREIGN KEY fk_summary_category,
    DROP INDEX idx_subcategory_id,
    DROP INDEX idx_category_id,
    DROP COLUMN subcategory_id,
    DROP COLUMN category_id;

DROP TABLE IF EXISTS subcategories;
DROP TABLE IF EXISTS categories;

