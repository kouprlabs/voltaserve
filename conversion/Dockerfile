FROM registry.suse.com/bci/golang:1.22 AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-conversion

FROM opensuse/leap:15.5 AS runner

COPY --from=builder /build/scripts ./scripts

RUN ./scripts/install-dependencies.sh
RUN ./scripts/install-fonts.sh

WORKDIR /app

COPY --from=builder /build/voltaserve-conversion ./voltaserve-conversion
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-conversion"]

EXPOSE 5001
