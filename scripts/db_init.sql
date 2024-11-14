CREATE DATABASE social2
CREATE EXTENSION IF NOT EXISTS citext;


-- Create migration: migrate create -ext sql -dir cmd/migrate/migrations -seq create_users_table

-- Make migration: migrate -database postgres://user:password@localhost/db_name?sslmode=disable -path ./cmd/api/migrate/migrations up
