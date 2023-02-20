FROM node:18-alpine as builder

WORKDIR /build

COPY . .

RUN npm install
RUN npm run pkg

FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV production

COPY --from=builder /build/pkg/voltaserve-idp-alpine ./voltaserve-idp
COPY --from=builder /build/templates ./templates
COPY --from=builder /build/config.yml ./config.yml

ENTRYPOINT ["./voltaserve-idp"]

EXPOSE 7000