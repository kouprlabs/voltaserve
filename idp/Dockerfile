FROM node:18-alpine

WORKDIR /app

COPY . .

RUN npm install --production

ENTRYPOINT ["/bin/sh", "-c", "npm run start"]

EXPOSE 7000