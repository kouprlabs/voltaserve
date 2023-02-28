FROM node:18-alpine

WORKDIR /app

COPY . .

RUN npm install

ENTRYPOINT ["/bin/sh", "-c", "npm run dev"]

EXPOSE 3000