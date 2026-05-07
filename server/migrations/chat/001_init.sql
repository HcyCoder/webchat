CREATE DATABASE IF NOT EXISTS webchat_chat;
USE webchat_chat;

CREATE TABLE IF NOT EXISTS messages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    chat_type VARCHAR(10) NOT NULL,
    from_user BIGINT NOT NULL,
    to_id BIGINT NOT NULL,
    msg_type VARCHAR(20) NOT NULL,
    content TEXT,
    is_recalled TINYINT DEFAULT 0,
    created_at BIGINT NOT NULL,
    INDEX idx_to_id_chat_type_created (to_id, chat_type, created_at)
);

CREATE TABLE IF NOT EXISTS conversations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    chat_type VARCHAR(10) NOT NULL,
    target_id BIGINT NOT NULL,
    last_msg_id BIGINT DEFAULT 0,
    unread_count INT DEFAULT 0,
    is_pinned TINYINT DEFAULT 0,
    is_muted TINYINT DEFAULT 0,
    updated_at BIGINT NOT NULL,
    UNIQUE KEY uk_user_conv (user_id, chat_type, target_id)
);

CREATE TABLE IF NOT EXISTS read_receipts (
    msg_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    read_at BIGINT NOT NULL,
    PRIMARY KEY (msg_id, user_id)
);
