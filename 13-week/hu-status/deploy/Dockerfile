# ---- Build stage ----
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm install --omit=dev

COPY src ./src

# ---- Run stage ----
FROM node:22-alpine

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/src ./src

EXPOSE 8080

CMD ["node", "src/index.js"]
