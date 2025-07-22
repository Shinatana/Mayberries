GOOSE_MIGRATION_DIR=/Users/konoko/Documents/golang/Mayberries/auth_service/internal/repo/migrate
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://myuser:secret@localhost:5432/mydb?sslmode=disable&search_path=public

.PHONY: migrate migrate_down create_migration

migrate:
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" up

migrate_down:
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" down

create_migration:
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) "$(GOOSE_DBSTRING)" create $(name) sql

.PHONY: migrate migrate_down create_migration