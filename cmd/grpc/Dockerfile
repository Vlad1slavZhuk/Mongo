FROM golang:alpine AS builder
WORKDIR /app
COPY ./cmd/grpc /app/cmd/grpc
COPY ./internal /app/internal
COPY ./api /app/api
COPY ./go.mod /app/go.mod
COPY ./vendor /app/vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o grpc /app/cmd/grpc/main.go


FROM alpine:latest
WORKDIR /gRPC
COPY --from=builder /app/grpc .
COPY --from=builder /app/internal/pkg/migrations ./migrations
CMD [ "./grpc" ]