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

kube-service:
	kubectl apply -f ./build/$n/$n-service.yaml

kube-deployment:
	kubectl apply -f ./build/$n/$n-deployment.yaml

kube-ingress:
	kubectl apply -f ./build/hello-sekai-shop-ingress.yaml

kube-create-configmap:
	kubectl create configmap app-env --from-file=./env/prod/.env

kube-delete-configmap:
	kubectl delete configmap $n

# Run all services
run-all:
	go run server/auth.go ./env/$(e)/.env.$(a) & \
	go run server/inventory.go ./env/$(e)/.env.$(a) & \
	go run server/item.go ./env/$(e)/.env.$(a) & \
	go run server/payment.go ./env/$(e)/.env.$(a) & \
	go run server/player.go ./env/$(e)/.env.$(a)

# Kill all running Go processes
kill-all:
	pkill -f "go run server/"