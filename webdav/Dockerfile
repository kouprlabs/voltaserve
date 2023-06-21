FROM registry.suse.com/bci/nodejs:18-4.3

WORKDIR /app

COPY src ./src
COPY .env .
COPY package.json .
COPY pnpm-lock.yaml .
COPY tsconfig.json .

RUN npm install -g corepack
RUN corepack enable

RUN pnpm install

ENTRYPOINT ["pnpm", "run", "start"]

EXPOSE 7000
