FROM golang:1.23-rc-alpine

RUN apk --no-cache add git

RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

WORKDIR /app

COPY go.mod go.sum .env ./

COPY api/ ./api/ 

COPY gen/ ./gen/

COPY cmd/server/ ./cmd/server/

COPY services/mqtt/ ./services/mqtt/

COPY db/ ./db/

COPY pkg/ ./pkg/

RUN go mod download

RUN go mod tidy

RUN go build -o server ./cmd/server/main.go

EXPOSE 8000

EXPOSE 8001

CMD ["./server"]