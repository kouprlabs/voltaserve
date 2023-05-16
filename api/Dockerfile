FROM golang:1.20-bullseye as builder

WORKDIR /build

COPY . .

ENV GOPROXY=direct

RUN go mod download
RUN go build -o voltaserve-api

FROM debian:bullseye AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-api ./voltaserve-api
COPY --from=builder /build/install-deps.sh ./install-deps.sh
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/templates ./templates

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8
ENV DEBIAN_FRONTEND=noninteractive

RUN ./install-deps.sh

ENTRYPOINT ["./voltaserve-api"]

EXPOSE 5000