sqlc:
	sqlc generate

migrateup:
	migrate -path backend/db/migrate -database "mysql://root:123@tcp(localhost:3306)/candles" -verbose up

migratedown:
	migrate -path backend/db/migrate -database "mysql://root:123@tcp(localhost:3306)/candles" -verbose down

.PHONY: sqlc