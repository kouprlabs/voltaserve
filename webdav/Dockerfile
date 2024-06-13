FROM node:22-alpine

WORKDIR /app

COPY src ./src
COPY .env .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .

RUN npm ci --legacy-peer-deps

ENTRYPOINT ["npm", "run", "start"]

EXPOSE 8082
