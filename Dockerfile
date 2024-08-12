FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags=viper_bind_struct -o prizer


FROM alpine:latest

COPY --from=build /app/prizer /prizer

ENV GIN_MODE=release

ENTRYPOINT ["/prizer"]