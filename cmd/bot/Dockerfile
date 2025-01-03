FROM golang:1.23 AS builder

WORKDIR /shortybot

COPY ./ ./

RUN go mod download
RUN go build -o shortybot ./cmd/bot/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /shortybot

COPY --from=builder /shortybot/shortybot ./

CMD ["./shortybot"]