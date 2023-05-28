sqlc:
	sqlc generate

migrateup:
	migrate -path backend/db/migrate -database "mysql://root:123@tcp(localhost:3306)/bitmoi" -verbose up

migratedown:
	migrate -path backend/db/migrate -database "mysql://root:123@tcp(localhost:3306)/bitmoi" -verbose down

migrateup1:
	migrate -path backend/db/migrate -database "mysql://root:123@tcp(localhost:3306)/bitmoi" -verbose up 1

migratedown1:
	migrate -path backend/db/migrate -database "mysql://root:123@tcp(localhost:3306)/bitmoi" -verbose down 1

mock:
	mockgen -package mockdb -destination backend/db/mock/store.go bitmoi/backend/db/sqlc Store

.PHONY: sqlc migrateup migratedown migrateup1 migratedown1 mock