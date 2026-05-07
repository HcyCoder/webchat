CREATE DATABASE IF NOT EXISTS webchat_group;
USE webchat_group;

CREATE TABLE IF NOT EXISTS grps (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    avatar VARCHAR(500) DEFAULT '',
    owner_id BIGINT NOT NULL,
    announcement VARCHAR(500) DEFAULT '',
    member_count INT DEFAULT 0,
    max_members INT DEFAULT 500,
    created_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS group_members (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    alias VARCHAR(50) DEFAULT '',
    is_muted TINYINT DEFAULT 0,
    joined_at BIGINT NOT NULL,
    UNIQUE KEY uk_group_user (group_id, user_id)
);
