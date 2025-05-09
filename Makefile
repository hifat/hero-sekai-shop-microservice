run:
	go run . ./env/$e/.env.$a

pb-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./moduels/$f/$mProto/$m.proto

seed:
	go run ./pkg/database/script/migration.go ./env/$e/.env.$a

db-up:
	docker compose -f docker-compose.db.yaml up -d

db-down:
	docker compose -f docker-compose.db.yaml down

kafka-up:
	docker compose -f docker-compose.kafka.yaml up -d

kafka-down:
	docker compose -f docker-compose.kafka.yaml down

d-build:
	docker build -f ./build/Dockerfile -t hellosekai-shop:latest .

d-push:
	docker push $u/hello-sekai-shop-$s