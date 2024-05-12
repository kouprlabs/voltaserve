FROM oven/bun:1-debian

RUN apt-get update
RUN apt-get install -y curl

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY bun.lockb .
COPY tsconfig.json .

RUN bun install --frozen-lockfile

ENTRYPOINT ["bun", "run", "start"]

EXPOSE 7000
