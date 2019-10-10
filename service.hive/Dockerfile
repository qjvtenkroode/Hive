# develop stage
FROM node:11.1-alpine AS develop-stage
WORKDIR /app
COPY package*.json ./
RUN yarn install
COPY . .

FROM node:latest AS build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./ .
RUN npm run build

FROM nginx AS production-stage
RUN mkdir /app
COPY --from=build-stage /app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf
