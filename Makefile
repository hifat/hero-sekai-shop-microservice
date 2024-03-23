run:
	go run . ./env/$e/.env.$a

pb-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./moduels/$f/$fProto/$f.proto

seed:
	go run ./pkg/database/script/migration.go ./env/$e/.env.$a

db-up:
	docker compose -f docker-compose.db.yaml up
