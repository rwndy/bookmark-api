DB_URL=postgres://postgres@localhost:5432/bookmark_db?sslmode=disable

run:
	go run cmd/api/main.go

migrate-up:
	migrate -database "$(DB_URL)" -path migrations up

migrate-down:
	migrate -database "$(DB_URL)" -path migrations down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)