gen: ./internal/app/proto/services.proto
	export PATH="$PATH:$(go env GOPATH)/bin"
	protoc --go_out=./internal/app/proto/gen/ --go-grpc_out=./internal/app/proto/gen/ ./internal/app/proto/services.proto