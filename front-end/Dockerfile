# Stage 1: Build the front-end
FROM node:18-alpine as build

WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Serve the built front-end with Nginx
FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
# Copy a custom nginx config if you have one (optional)
# COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]