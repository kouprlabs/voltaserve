FROM opensuse/leap:15.5 as builder

WORKDIR /build

COPY . .

RUN zypper install -y go1.20 git

RUN go mod download
RUN go build -o voltaserve-conversion

FROM opensuse/leap:15.5 AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-conversion ./voltaserve-conversion
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-conversion"]

EXPOSE 5001