FROM node:22-bookworm

RUN apt-get update
RUN apt-get install -y curl

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .

RUN npm ci --legacy-peer-deps

ENTRYPOINT ["npm", "run", "start"]

EXPOSE 8081
