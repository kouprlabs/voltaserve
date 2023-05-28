FROM golang:1.20-alpine as builder

WORKDIR /build

COPY . .

ENV GOPROXY=direct

RUN go mod download
RUN go build -o voltaserve-api

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-api ./voltaserve-api
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/templates ./templates

ENTRYPOINT ["./voltaserve-api"]

EXPOSE 5000