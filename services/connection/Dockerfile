FROM golang:1.23-rc-alpine

RUN apk --no-cache add git

WORKDIR /app

COPY go.mod go.sum .env ./

COPY services/connection/ ./services/connection/

COPY gen/ ./gen/

COPY database/ ./database/

COPY pkg/ ./pkg/

RUN go mod download

RUN go mod tidy

RUN go build -o connection ./services/connection/main.go

EXPOSE 8004

CMD ["./connection"]