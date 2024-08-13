FROM golang:1.23-rc-alpine

RUN apk --no-cache add git

RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

WORKDIR /app

COPY go.mod go.sum .env ./

COPY services/chat/ ./services/chat/

COPY gen/ ./gen/

COPY db/ ./db/

COPY pkg/ ./pkg/

RUN go mod download

RUN go mod tidy

RUN go build -o chat ./services/chat/main.go

EXPOSE 8002

CMD ["./chat"]