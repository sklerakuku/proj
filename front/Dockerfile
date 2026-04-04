# Этап 1: Сборка статики
FROM node:20-alpine AS build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Этап 2: Раздача через Nginx
FROM nginx:stable-alpine as production-stage
# Копируем билд из первого этапа в папку nginx
COPY --from=build-stage /app/dist /usr/share/nginx/html
# Копируем конфиг (раскомментируй эту строку!)
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]