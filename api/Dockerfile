FROM redhat/ubi9 as builder

WORKDIR /build

COPY . .

RUN dnf install -y golang

RUN go mod download
RUN go build -o voltaserve-api

FROM redhat/ubi9 AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-api ./voltaserve-api
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/templates ./templates

ENTRYPOINT ["./voltaserve-api"]

EXPOSE 5000