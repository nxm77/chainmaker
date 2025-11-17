-- 创建数据库
CREATE DATABASE IF NOT EXISTS chainmaker_dquery CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS chainmaker_explorer CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS contract_compile CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建只读用户并授予所有数据库的只读权限
CREATE USER IF NOT EXISTS 'readonly'@'%' IDENTIFIED BY 'readonly';
GRANT SELECT ON *.* TO 'readonly'@'%';  -- 赋予所有数据库的只读权限

-- 创建 chainmaker 用户并授予所有权限
CREATE USER IF NOT EXISTS 'chainmaker'@'%' IDENTIFIED BY 'chainmaker';
GRANT ALL PRIVILEGES ON *.* TO 'chainmaker'@'%';  -- 赋予所有权限
FLUSH PRIVILEGES;  -- 刷新权限
