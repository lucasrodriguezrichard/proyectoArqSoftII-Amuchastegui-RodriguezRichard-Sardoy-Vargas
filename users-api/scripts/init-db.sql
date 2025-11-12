-- Initialization script for users_db
-- This is optional since GORM auto-migrates, but useful for reference

USE users_db;

-- The users table will be auto-created by GORM
-- This is just for documentation

-- CREATE TABLE IF NOT EXISTS users (
--     id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
--     username VARCHAR(32) UNIQUE NOT NULL,
--     email VARCHAR(191) UNIQUE NOT NULL,
--     first_name VARCHAR(100) NOT NULL,
--     last_name VARCHAR(100) NOT NULL,
--     password_hash VARCHAR(255) NOT NULL,
--     role VARCHAR(20) DEFAULT 'user' NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
--     INDEX idx_username (username),
--     INDEX idx_email (email)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SELECT 'users_db initialized' AS message;
