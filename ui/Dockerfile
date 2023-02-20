FROM node:18-alpine

WORKDIR /app

COPY . .

ARG API_URL
ARG IDP_URL

ENV API_URL=$API_URL
ENV IDP_URL=$IDP_URL

ENV NODE_ENV production

RUN npm install

ENTRYPOINT ["/bin/sh", "-c", "npm run dev"]

EXPOSE 5173