FROM golang:1.22 AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-api

FROM ubuntu:noble AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-api ./voltaserve-api
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/templates ./templates

ENTRYPOINT ["./voltaserve-api"]

EXPOSE 5000
