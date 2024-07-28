FROM golang:1.23-rc-alpine

RUN apk --no-cache add git

WORKDIR /app

COPY ../../go.mod ../../go.sum ./

COPY . .

RUN go mod tidy

RUN go build -o chat .

EXPOSE 8002

CMD ["./chat"]