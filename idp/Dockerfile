FROM node:20-alpine

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .

RUN npm ci --legacy-peer-deps

ENTRYPOINT ["npm", "run", "start:node"]

EXPOSE 8081
