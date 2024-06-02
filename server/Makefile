MIGRATION_PATH = ./db/migrations

postgres:
	docker compose up

migrateup:
	goose -dir $(MIGRATION_PATH) postgres "user=postgres dbname=boards sslmode=disable password=postgres host=localhost port=5433" up

migratedown:
	goose -dir $(MIGRATION_PATH) postgres "user=root dbname=boards sslmode=disable password=postgres host=localhost port=5433" down	
