ifeq ($(OS),Windows_NT)
	ifneq ($(ls backend/gapi/pb),)
		DELETE_COMMAND=cd backend/gapi/pb&&del *pb.go
	else
		DELETE_COMMAND=echo "already empty"
	endif
else
	DELETE_COMMAND=rm backend/gapi/pb/*pb.go
endif

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

delete:
	$(DELETE_COMMAND)

proto: delete
	protoc \
	--proto_path=backend/gapi/proto --go_out=backend/gapi/pb \
	--go_opt=paths=source_relative --go-grpc_out=backend/gapi/pb \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=backend/gapi/pb --grpc-gateway_opt=paths=source_relative \
	--validate_out="lang=go:backend/gapi/pb" --validate_opt=paths=source_relative \
	backend/gapi/proto/*.proto

checkos:
	echo $(OS)
.PHONY: sqlc migrateup migratedown migrateup1 migratedown1 mock proto checkos