FROM golang:alpine AS builder
WORKDIR /app
COPY ./cmd/http /app/cmd/http
COPY ./internal /app/internal
COPY ./api /app/api
COPY ./go.mod /app/go.mod
COPY ./vendor /app/vendor
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o http /app/cmd/http/main.go


FROM alpine:latest
WORKDIR /HTTP
COPY --from=builder /app/http .
CMD [ "./http" ]