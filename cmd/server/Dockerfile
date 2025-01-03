FROM golang:1.23 AS builder

WORKDIR /shortyserver

COPY ./ ./

RUN go mod download
RUN go build -o shortyserver ./cmd/server/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /shortyserver

COPY --from=builder /shortyserver/shortyserver ./

EXPOSE 8080

CMD ["./shortyserver"]