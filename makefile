create_migration:
# make create_migration name=name_your_migration_without_spaces
	migrate create -ext sql -dir db/migrations -seq ${name}
migrate:
# make migrate password=postgres_password
	migrate -database 'postgres://postgres:${password}@localhost:5430/user_service?sslmode=disable' -path ./db/migrations up
fmt:
	go fmt ./...
local:
	go build -o . cmd/main.go
	./main --use_db_config