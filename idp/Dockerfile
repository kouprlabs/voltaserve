FROM opensuse/leap:15.5

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY pnpm-lock.yaml .
COPY tsconfig.json .

RUN zypper install -y nodejs18

RUN npm install -g corepack
RUN corepack enable

RUN pnpm install

ENTRYPOINT ["pnpm", "run", "start"]

EXPOSE 7000
