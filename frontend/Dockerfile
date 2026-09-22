FROM node:20-alpine AS builder

WORKDIR /app
# Use wildcard to copy package.json and package-lock.json (if it exists)
COPY package*.json ./
# Use install instead of ci since we might not have a lockfile yet
RUN npm install

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# React router support
RUN echo 'server { listen 80; location / { root /usr/share/nginx/html; index index.html index.htm; try_files $uri $uri/ /index.html; } }' > /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
