SERVER_PATH = server/main.go
SERVER_PORT = 50051/tcp
WEBSERVER_PORT = 8003/tcp
CSV_PATH = http://localhost:8003/products.csv


up:
	docker-compose up --scale app=3

build:
	docker-compose build

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		grpcserver/grpc_server.proto

local_server:
	go run $(SERVER_PATH)&

local_web_server:
	cd test && ./web-server.sh&

kill_local_webserver:
	fuser -k $(WEBSERVER_PORT)

kill_local_server:
	fuser -k $(SERVER_PORT)
	fuser -k $(WEBSERVER_PORT)

local_client:
	go run client/main.go  $(CSV_PATH)

run_local_client: local_web_server local_client kill_local_webserver

local_test: local_server run_local_client kill_local_server

