FROM golang:1.23.2-alpine AS builder

RUN apk --no-cache add git

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go test ./... -v
RUN go build -o /card-validator ./cmd

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /card-validator .

EXPOSE 9000
CMD ["./card-validator"]