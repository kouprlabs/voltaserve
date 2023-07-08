FROM registry.suse.com/bci/nodejs:18

WORKDIR /app

COPY src ./src
COPY public ./public
COPY index.html .
COPY postcss.config.cjs .
COPY tailwind.config.cjs .
COPY vite.config.ts .
COPY .env .
COPY package.json .
COPY pnpm-lock.yaml .
COPY tsconfig.json .
COPY tsconfig.node.json .

RUN npm install -g corepack
RUN corepack enable

RUN pnpm install

ENTRYPOINT ["pnpm", "run", "dev"]

EXPOSE 3000
