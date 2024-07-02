FROM oven/bun:1-alpine AS ui-builder

WORKDIR /build

COPY src ./src
COPY public ./public
COPY index.html .
COPY postcss.config.cjs .
COPY tailwind.config.cjs .
COPY vite.config.ts .
COPY package.json .
COPY bun.lockb .
COPY tsconfig.json .
COPY tsconfig.node.json .

RUN bun install --frozen-lockfile --production
RUN bun run build

FROM golang:1.22-alpine AS go-builder

WORKDIR /build

COPY --from=ui-builder /build/dist ./dist
COPY config ./config
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build -o voltaserve-ui

FROM alpine:3.20 AS runner

WORKDIR /app

COPY --from=go-builder /build/voltaserve-ui ./voltaserve-ui
COPY .env .

ENTRYPOINT ["./voltaserve-ui"]

EXPOSE 3000
