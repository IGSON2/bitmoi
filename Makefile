ifneq ($(ls backend/gapi/pb),)
	DELETE_COMMAND=cd backend/gapi/pb&&del *pb.go
else
	DELETE_COMMAND=$(echo "already empty")
endif

sqlc:
	sqlc generate
	make mock

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
	mockgen -package mocktask -destination backend/worker/mock/distributor.go bitmoi/backend/worker TaskDistributor

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

installgrpcweb:
	npm install grpc-tools grpc-web

reactproto:
	npx grpc_tools_node_protoc \
	--js_out=import_style=commonjs,binary:frontend/src/component/pb \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:frontend/src/component/pb \
	--proto_path=backend/gapi/proto \
	backend/gapi/proto/*.proto

rmi:
	docker compose down && docker rmi bitmoi_api

test:
	go test -v -cover -short ./backend/...

benchmark:
	go-wrk -c 80 -d 5 -H Content-Type:application/json -M GET http://43.202.77.76:5000/practice

swag:
	rm -rf ./frontend/server/docs
	swag init --output ./frontend/server/docs

client:
	cd ./frontend/server&&go run main.go

.PHONY: sqlc migrateup migratedown migrateup1 migratedown1 mock proto installgrpcweb reactproto rmi test benchmark swag client