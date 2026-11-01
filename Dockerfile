FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download || true

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/mesh-sentry ./cmd/daemon

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/mesh-sentry .

EXPOSE 9090
ENTRYPOINT ["./mesh-sentry"]
