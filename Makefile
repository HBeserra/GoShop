start_services:
	docker-compose -f ./deployments/docker-compose.dev.yaml up -d

stop_services:
	docker-compose -f ./deployments/docker-compose.dev.yaml down

run_sqlc:
	sqlc generate

migrate_up:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goshop?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goshop?sslmode=disable" -verbose down

migrate_drop:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goshop?sslmode=disable" -verbose drop


.PHONY: start_services stop_services run_sqlc migrate_up migrate_down migrate_drop
