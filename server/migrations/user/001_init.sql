CREATE DATABASE IF NOT EXISTS webchat_user;
USE webchat_user;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    avatar VARCHAR(500) DEFAULT '',
    gender TINYINT DEFAULT 0,
    region VARCHAR(100) DEFAULT '',
    signature VARCHAR(200) DEFAULT '',
    created_at BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS contacts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    contact_id BIGINT NOT NULL,
    remark VARCHAR(50) DEFAULT '',
    tag VARCHAR(50) DEFAULT '',
    is_blocked TINYINT DEFAULT 0,
    added_at BIGINT NOT NULL,
    UNIQUE KEY uk_user_contact (user_id, contact_id)
);

CREATE TABLE IF NOT EXISTS friend_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    from_user BIGINT NOT NULL,
    to_user BIGINT NOT NULL,
    message VARCHAR(100) DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at BIGINT NOT NULL
);
