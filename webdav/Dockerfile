FROM node:18-alpine as builder

WORKDIR /build

COPY . .

RUN yarn run pkg

FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV production

COPY --from=builder /build/pkg/voltaserve-webdav-alpine ./voltaserve-webdav
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-webdav"]

EXPOSE 6000
