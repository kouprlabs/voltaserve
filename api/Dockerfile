FROM opensuse/leap:15.5 as builder

WORKDIR /build

COPY . .

RUN zypper install -y go1.20 git

RUN go mod download
RUN go build -o voltaserve-api

FROM opensuse/leap:15.5 AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-api ./voltaserve-api
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/templates ./templates

ENTRYPOINT ["./voltaserve-api"]

EXPOSE 5000