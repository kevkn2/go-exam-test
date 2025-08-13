.PHONY: migrate-create migrate-up migrate-down

MIGRATION_DIR = ./migrations # Directory where migration files will be stored
DATABASE_URL = "postgres://exam_test:amb4tukAm@localhost:5433/exam_test_db?sslmode=disable" # Replace with your actual database URL

# Create a new migration file
migrate-create:
	@echo "Creating new migration: $(name)..."
	migrate create -ext sql -dir $(MIGRATION_DIR) $(name)

# Apply pending migrations
migrate-up:
	@echo "Applying migrations..."
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) up

# Rollback the last migration
migrate-down:
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATION_DIR) -database $(DATABASE_URL) down 1