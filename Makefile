include .env
export

run:
	go run cmd/api/main.go

migrate-up:
	migrate -database "$(DB_DSN)" -path migrations up

migrate-down:
	migrate -database "$(DB_DSN)" -path migrations down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)