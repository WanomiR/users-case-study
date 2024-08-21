FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init userservice && go mod tidy

RUN go build -o main ./cmd/api

FROM alpine:latest

COPY --from=builder /app/main /main
COPY --from=builder /app/docs /docs

EXPOSE 8888

CMD ["/main"]
