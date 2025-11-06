# ===================================
# Stage 1: Build stage
# ===================================
FROM node:20-alpine AS builder

WORKDIR /app

# Accept build arguments
ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# Install dependencies first (better layer caching)
COPY package.json package-lock.json ./
RUN npm ci

# Copy source code
COPY . .

# Build the application (vite build only, skip type checking for Docker)
ENV NODE_ENV=production
RUN npx vite build

# ===================================
# Stage 2: Production stage
# ===================================
FROM nginx:alpine

# Copy custom nginx config for Vue Router support
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:80/ || exit 1

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
