FROM oven/bun:1-alpine

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY bun.lockb .
COPY tsconfig.json .

RUN bun install --frozen-lockfile --production

ENTRYPOINT ["bun", "run", "start"]

EXPOSE 8081
