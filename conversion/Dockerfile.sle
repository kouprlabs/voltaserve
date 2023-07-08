FROM registry.suse.com/bci/golang:1.20 AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-conversion

FROM registry.suse.com/bci/bci-micro:15.5 AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-conversion ./voltaserve-conversion
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-conversion"]

EXPOSE 5001
