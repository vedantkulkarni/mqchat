gen: ./internal/app/proto/services.proto
	export PATH="$PATH:$(go env GOPATH)/bin"
	protoc --go-grpc_out=./internal/app/proto/gen/ ./internal/app/proto/services.proto