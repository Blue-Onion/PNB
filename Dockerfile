# syntax=docker/dockerfile:1

FROM golang:1.27.2-bookworm AS backend
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -o /out/txshield-api ./cmd

FROM oven/bun:1.3.14 AS bun

FROM node:22-bookworm-slim AS frontend
COPY --from=bun /usr/local/bin/bun /usr/local/bin/bun
WORKDIR /app
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile
COPY frontend/ ./
# Next.js writes the rewrite destination into its production build.
ENV BACKEND_URL=http://127.0.0.1:8080
ENV NEXT_TELEMETRY_DISABLED=1
RUN bun run build

FROM node:22-bookworm-slim AS runtime
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
COPY --from=backend /out/txshield-api /usr/local/bin/txshield-api
COPY --from=frontend --chown=node:node /app/.next ./.next
COPY --from=frontend --chown=node:node /app/node_modules ./node_modules
COPY --from=frontend /app/package.json /app/next.config.ts ./
COPY --from=frontend /app/public ./public
COPY docker-entrypoint.sh /usr/local/bin/txshield-start
USER node
EXPOSE 3000
ENTRYPOINT ["/bin/bash", "/usr/local/bin/txshield-start"]
