FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-api

FROM golang:1.22-alpine AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-api ./voltaserve-api
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/templates ./templates

ENTRYPOINT ["./voltaserve-api"]

EXPOSE 8080
