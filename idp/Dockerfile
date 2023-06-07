FROM redhat/ubi9

WORKDIR /app

COPY src .
COPY templates .
COPY .env .
COPY package.json .
COPY pnpm-lock.yaml .
COPY tsconfig.json .

RUN dnf module -y enable nodejs:18
RUN dnf module -y install nodejs:18/common

RUN npm install -g corepack
RUN corepack enable

RUN pnpm install

ENTRYPOINT ["pnpm", "run", "start"]

EXPOSE 7000
