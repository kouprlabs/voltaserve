FROM oven/bun:1

WORKDIR /app

COPY src ./src
COPY .env .
COPY package.json .
COPY bun.lockb .
COPY tsconfig.json .

RUN bun install --frozen-lockfile

ENTRYPOINT ["bun", "run", "dev"]

EXPOSE 6000
