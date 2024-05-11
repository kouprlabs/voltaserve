FROM oven/bun:1

WORKDIR /app

COPY src ./src
COPY public ./public
COPY index.html .
COPY postcss.config.cjs .
COPY tailwind.config.cjs .
COPY vite.config.ts .
COPY .env .
COPY package.json .
COPY bun.lockb .
COPY tsconfig.json .
COPY tsconfig.node.json .

RUN bun install --frozen-lockfile

ENTRYPOINT ["bun", "run", "dev"]

EXPOSE 3000
