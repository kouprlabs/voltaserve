FROM opensuse/leap:15.5 as builder

WORKDIR /build

COPY . .

RUN zypper install -y go1.20

RUN go mod download
RUN go build -o voltaserve-tools

FROM opensuse/leap:15.5 AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-tools ./voltaserve-tools
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-tools"]

EXPOSE 6001