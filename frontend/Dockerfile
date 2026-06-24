FROM node:22-alpine AS base

FROM base AS installer
WORKDIR /app

# Копируем package.json (lock-файл от pnpm копировать НЕ нужно)
COPY package.json ./

# Ставим зависимости через чистый npm
RUN npm install

COPY . .
RUN npm run build

FROM base AS runner
WORKDIR /app
ENV NODE_ENV=production

COPY --from=installer /app/public ./public
COPY --from=installer /app/.next/standalone ./
COPY --from=installer /app/.next/static ./.next/static

EXPOSE 3000
CMD ["node", "server.js"]
