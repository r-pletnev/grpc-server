SERVER_PATH = server/main.go
SERVER_PORT = 50051/tcp
WEBSERVER_PORT = 8003/tcp
CSV_PATH = http://localhost:8003/products.csv



proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		grpcserver/grpc_server.proto

srv:
	go run $(SERVER_PATH)&

web_srv:
	cd test && ./web-server.sh&

kill_webserver:
	fuser -k $(WEBSERVER_PORT)

kill_server:
	fuser -k $(SERVER_PORT)
	fuser -k $(WEBSERVER_PORT)

clt:
	go run client/main.go  $(CSV_PATH)

run_client: web_srv clt kill_webserver


test: srv run_client kill_server

