@echo off
migrate -path ./schema -database 'postgres://postgres:postgres@192.168.0.104:5432/postgres?sslmode=disable' up