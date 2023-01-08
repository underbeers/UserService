create_migration:
# make create_migration name=name_your_migration_without_spaces
	migrate create -ext sql -dir db/migrations -seq ${name}
migrate:
# make migrate password=postgres_password host=localhost port=5420 mode=up/down
	migrate -database 'postgres://postgres:${password}@${host}:${port}/user_service?sslmode=disable' -path ./db/migrations ${mode}
fmt:
	go fmt ./...
db_win:
	docker run -d --name=pgsql -p 5420:5432 -e POSTGRES_PASSWORD='localpassword' -v C:\Program_Files\PostgreSQL\14\data:/var/lib/postgresql/data postgres
db_unix:
	docker run -d --name=pgsql -p 5420:5432 -e POSTGRES_PASSWORD='localpassword' -v /var/lib/pgsql/data:/var/lib/pgsql/data postgres
local:
	go build -o . cmd/main.go
	./main --use_db_config
build_image:
	docker build -t rodmul/pl_user_service .
run:
	docker run -d -p 6001:6001 -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
	-e POSTGRES_HOST=$POSTGRES_HOST -e POSTGRES_USER=$POSTGRES_USER \
	-e POSTGRES_PORT=$POSTGRES_PORT -e POSTGRES_DB_NAME=$POSTGRES_DB_NAME \
	--name user_service_container rodmul/pl_user_service