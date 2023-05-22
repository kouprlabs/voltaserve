FROM node:18-alpine as builder

WORKDIR /build

COPY . .

RUN npm install
RUN npm run pkg

FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV production

COPY --from=builder /build/pkg/voltaserve-webdav-alpine ./voltaserve-webdav

ENTRYPOINT ["./voltaserve-webdav"]

EXPOSE 6000