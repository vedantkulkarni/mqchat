FROM golang:1.23-rc-alpine

RUN apk --no-cache add git

WORKDIR /app

COPY go.mod go.sum .env ./

COPY services/user/ ./services/user/

COPY gen/ ./gen/

COPY database/ ./database/

COPY pkg/ ./pkg/

RUN go mod download

RUN go mod tidy

RUN go build -o user ./services/user/main.go

EXPOSE 8003

CMD ["./user"]