FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-mosaic

FROM golang:1.22-alpine AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-mosaic ./voltaserve-mosaic
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-mosaic"]

EXPOSE 8085
