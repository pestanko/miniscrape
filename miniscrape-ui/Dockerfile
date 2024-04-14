
# syntax=docker/dockerfile:1
# https://docs.docker.com/develop/develop-images/multistage-build/

ARG NODE_VERSION=20-slim

# BASE
FROM node:$NODE_VERSION AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /app

COPY ../package.json ../pnpm-lock.yaml ./

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

# BUILD
FROM base AS build

WORKDIR /app

COPY . .

COPY --chown=node:node --from=base /app/node_modules ./node_modules

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

# RUNTIME
FROM node:$NODE_VERSION

WORKDIR /app

COPY --chown=node:node --from=build /app/package.json /app/pnpm-lock.yaml ./
COPY --chown=node:node --from=build /app/node_modules ./node_modules
COPY --chown=node:node --from=build /app/dist ./dist

EXPOSE 3000

USER node

CMD ["node", "dist/index.js"]
