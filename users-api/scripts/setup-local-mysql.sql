-- Script para configurar MySQL local para users-api
-- Ejecutar con: mysql -u root -p < scripts/setup-local-mysql.sql

-- Crear base de datos
CREATE DATABASE IF NOT EXISTS users_db;

-- Crear usuario y otorgar permisos
CREATE USER IF NOT EXISTS 'restaurant_user'@'localhost' IDENTIFIED BY 'restaurant_pass';
GRANT ALL PRIVILEGES ON users_db.* TO 'restaurant_user'@'localhost';
FLUSH PRIVILEGES;

-- Usar la base de datos
USE users_db;

-- Verificar
SELECT 'users_db configurada correctamente!' AS message;
