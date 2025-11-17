-- 创建数据库
DO $$ 
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'chainmaker_dquery') THEN
      CREATE DATABASE chainmaker_dquery ENCODING 'UTF8';
   END IF;
END $$;

DO $$ 
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'chainmaker_explorer') THEN
      CREATE DATABASE chainmaker_explorer ENCODING 'UTF8';
   END IF;
END $$;

DO $$ 
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'contract_compile') THEN
      CREATE DATABASE contract_compile ENCODING 'UTF8';
   END IF;
END $$;

-- 创建只读用户并授予权限
DO $$ 
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'readonly') THEN
      CREATE USER readonly WITH PASSWORD 'readonly';
   END IF;
END $$;

-- 授予 readonly 用户对所有数据库的只读权限
GRANT CONNECT ON DATABASE template1 TO readonly;
GRANT CONNECT ON DATABASE template0 TO readonly;

-- 授予 readonly 用户对所有数据库中所有 schema 的只读权限
GRANT USAGE ON SCHEMA public TO readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly;

-- 自动为新表授予只读权限
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO readonly;

-- 创建 chainmaker 用户并授予所有权限
DO $$ 
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'chainmaker') THEN
      CREATE USER chainmaker WITH PASSWORD 'chainmaker';
   END IF;
END $$;

-- 授予 chainmaker 用户对所有数据库的所有权限
GRANT CONNECT ON DATABASE template1 TO chainmaker;
GRANT CONNECT ON DATABASE template0 TO chainmaker;

-- 授予 chainmaker 用户对所有数据库中所有 schema 的所有权限
GRANT ALL PRIVILEGES ON DATABASE chainmaker_dquery TO chainmaker;
GRANT ALL PRIVILEGES ON DATABASE chainmaker_explorer TO chainmaker;
GRANT ALL PRIVILEGES ON DATABASE contract_compile TO chainmaker;
GRANT ALL PRIVILEGES ON DATABASE template1 TO chainmaker;
GRANT ALL PRIVILEGES ON DATABASE template0 TO chainmaker;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO chainmaker;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO chainmaker;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO chainmaker;

-- 为新创建的表、序列、函数自动授予所有权限
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO chainmaker;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO chainmaker;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON FUNCTIONS TO chainmaker;
