# ==========================================================
# FRONTEND (DEVELOPMENT): Next.js dev server
# ==========================================================
FROM node:22-alpine

WORKDIR /app

ENV NEXT_TELEMETRY_DISABLED=1

# Dependency di-install saat build image. Saat container jalan,
# ./frontend di-mount ke /app, sedangkan node_modules memakai volume
# terpisah (lihat docker-compose.yml).
COPY package.json ./
RUN npm install

COPY . .

EXPOSE 3000
CMD ["npm", "run", "dev"]
