FROM registry.suse.com/bci/golang:1.20 as builder

WORKDIR /build

COPY . .

RUN zypper install -y go1.20

RUN go mod download
RUN go build -o voltaserve-conversion

FROM registry.suse.com/bci/bci-base:15.4 AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-conversion ./voltaserve-conversion
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-conversion"]

EXPOSE 5001