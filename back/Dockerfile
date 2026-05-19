FROM dockerhub.timeweb.cloud/library/golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/api/main.go

FROM dockerhub.timeweb.cloud/library/alpine:latest

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
