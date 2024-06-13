FROM node:22-alpine AS ui-builder

WORKDIR /build

COPY src ./src
COPY public ./public
COPY index.html .
COPY postcss.config.cjs .
COPY tailwind.config.cjs .
COPY vite.config.ts .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .
COPY tsconfig.node.json .

RUN npm ci --legacy-peer-deps
RUN npm run build

FROM golang:1.22-alpine AS go-builder

WORKDIR /build

COPY --from=ui-builder /build/dist ./dist
COPY config ./config
COPY main.go .
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build -o voltaserve-ui

FROM ubuntu:noble AS runner

RUN apt-get update
RUN apt-get install -y curl

WORKDIR /app

COPY --from=go-builder /build/voltaserve-ui ./voltaserve-ui
COPY .env .

ENTRYPOINT ["./voltaserve-ui"]

EXPOSE 3000
