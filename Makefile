gen : genUser genChat genRoom

genUser: ./proto/user.proto
	protoc --go_out=./gen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./gen/  \
	--go-grpc_opt=paths=source_relative \
	./proto/user.proto

genChat: ./proto/chat.proto
	protoc --go_out=./gen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./gen/  \
	--go-grpc_opt=paths=source_relative \
	./proto/chat.proto

genRoom: ./proto/room.proto
	protoc --go_out=./gen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./gen/  \
	--go-grpc_opt=paths=source_relative \
	./proto/room.proto
run:
	go run ./cmd/server/main.go