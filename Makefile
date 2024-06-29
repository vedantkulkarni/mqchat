gen : genUser genChat genConnection

genUser: ./proto/user.proto
	protoc --go_out=./internal/app/protogen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./internal/app/protogen/  \
	--go-grpc_opt=paths=source_relative \
	./proto/user.proto

genChat: ./proto/chat.proto
	protoc --go_out=./internal/app/protogen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./internal/app/protogen/  \
	--go-grpc_opt=paths=source_relative \
	./proto/chat.proto

genConnection: ./proto/connection.proto
	protoc --go_out=./internal/app/protogen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./internal/app/protogen/  \
	--go-grpc_opt=paths=source_relative \
	./proto/connection.proto
run:
	go run ./cmd/server/main.go