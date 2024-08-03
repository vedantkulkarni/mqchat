FROM golang:1.23-rc-alpine

RUN apk --no-cache add git

WORKDIR /app

COPY go.mod go.sum .env ./

COPY services/rooms/ ./services/rooms/

COPY gen/ ./gen/

COPY db/ ./db/

COPY pkg/ ./pkg/

RUN go mod download

RUN go mod tidy

RUN go build -o rooms ./services/rooms/main.go

EXPOSE 8004

CMD ["./rooms"]