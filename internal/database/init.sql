-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Run migrations
\i /docker-entrypoint-initdb.d/migrations/001_create_admin_user.sql 