FROM redhat/ubi9

WORKDIR /app

COPY src .
COPY public .
COPY postcss.config.cjs .
COPY tailwind.config.cjs .
COPY vite.config.js .
COPY .env .
COPY package.json .
COPY pnpm-lock.yaml .
COPY tsconfig.json .
COPY tsconfig.node.json .

RUN dnf module -y enable nodejs:18
RUN dnf module -y install nodejs:18/common

RUN npm install -g corepack
RUN corepack enable

RUN pnpm install

ENTRYPOINT ["/bin/sh", "-c", "pnpm run dev"]

EXPOSE 3000