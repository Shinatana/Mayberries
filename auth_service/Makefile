Migrations_Dir := "db/migrations"

.PHONY: migrate-generate
migrate-generate:
   $(GOPATH)/bin/goose -dir $(MIGRATIONS_DIR) create $(name) sql
