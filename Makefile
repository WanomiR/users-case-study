include .env

all: swagger_generate build

build:
	@docker compose up --build --force-recreate

swagger_generate:
	@cd userservice && swag init -g ./cmd/api/main.go && cd ..

migration_create:
	@migrate create -ext sql -dir postgres/migration/ -seq init_mg

migration_up:
	@migrate -path postgres/migration/ -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migration_down:
	@migrate -path postgres/migration/ -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

migration_fix:
	@migrate -path postgres/migration/ -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" force $(VERSION)

postgres_dump:
	@docker exec -i $(POSTGRES_CONT_NAME) /bin/sh -c "PGPASSWORD=$(DB_PASSWORD) pg_dump --username $(DB_USER) $(DB_NAME)" > ./postgres/create_tables.sql