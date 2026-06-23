# ──────────────────────────────────────────────
# Stage 1: Build the Quasar SPA
# ──────────────────────────────────────────────
FROM node:22-alpine AS builder

WORKDIR /app

# Install dependencies first (layer cache optimisation)
COPY package.json package-lock.json ./
RUN npm ci --ignore-scripts

# Copy source and build
COPY . .

# Build-time env vars are injected via --build-arg in CI
ARG VITE_USER_SERVICE
ARG VITE_AUTHOR_SERVICE
ARG VITE_BLOG_SERVICE
ARG VITE_GOOGLE_CLIENT_ID

RUN npx quasar build

# ──────────────────────────────────────────────
# Stage 2: Serve with Nginx
# ──────────────────────────────────────────────
FROM nginx:1.27-alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built SPA assets
COPY --from=builder /app/dist/spa /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
