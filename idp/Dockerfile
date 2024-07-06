FROM node:20-alpine AS base_amd64

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .

RUN npm i --legacy-peer-deps

FROM oven/bun:1-alpine AS base_arm64

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY bun.lockb .
COPY tsconfig.json .

RUN bun install --frozen-lockfile --production

FROM --platform=$TARGETPLATFORM base_${TARGETPLATFORM#*/} AS final

ENTRYPOINT if [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
      bun run start; \
    else \
      npm run start:node; \
    fi

EXPOSE 8081