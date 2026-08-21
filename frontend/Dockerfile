# Stage 1: install deps
FROM node:20-alpine AS deps
RUN corepack enable && corepack prepare pnpm@10.31.0 --activate
WORKDIR /app
# HORIZON_KEY: npm auth token for the private @squantumengine/* scope.
# Pass it from compose `build.args` (which interpolates from a .env file
# next to compose.yml — see nextjs/.env-example). Token is written to a
# transient .npmrc that is removed after install so it never lands in
# subsequent stages or the runner image.
ARG HORIZON_KEY
COPY package.json pnpm-lock.yaml ./
RUN printf '//registry.npmjs.org/:_authToken=%s\n@squantumengine:registry=https://registry.npmjs.org/\n' "$HORIZON_KEY" > .npmrc && \
    pnpm install --frozen-lockfile && \
    rm -f .npmrc

# Stage 2: build
FROM node:20-alpine AS builder
RUN corepack enable && corepack prepare pnpm@10.31.0 --activate
WORKDIR /app
# NEXT_PUBLIC_* vars must be present as ENV at build time — Next.js bakes them in during `next build`.
# ARG alone is not enough; we re-declare each var and export it as ENV before running the build.
ARG NEXT_PUBLIC_BACKEND_URL=https://api.example.com
ENV NEXT_PUBLIC_BACKEND_URL=${NEXT_PUBLIC_BACKEND_URL}
ARG NEXT_PUBLIC_RENDERING_MODE=ssr
ENV NEXT_PUBLIC_RENDERING_MODE=${NEXT_PUBLIC_RENDERING_MODE}
ARG NEXT_PUBLIC_APP_NAME=MyApp
ENV NEXT_PUBLIC_APP_NAME=${NEXT_PUBLIC_APP_NAME}
# boilerplate:axis=mocking option=msw START
ARG NEXT_PUBLIC_API_MOCKING=disabled
ENV NEXT_PUBLIC_API_MOCKING=${NEXT_PUBLIC_API_MOCKING}
# boilerplate:axis=mocking option=msw END
# boilerplate:axis=rendering option=csr START
ARG NEXT_PUBLIC_ENCRYPTION_KEY=dev-only-encryption-key-do-not-use-in-prod-32+
ENV NEXT_PUBLIC_ENCRYPTION_KEY=${NEXT_PUBLIC_ENCRYPTION_KEY}
# boilerplate:axis=rendering option=csr END
# boilerplate:axis=auth option=jwt-refresh START
ARG NEXT_PUBLIC_ENABLE_REFRESH_TOKEN=false
ENV NEXT_PUBLIC_ENABLE_REFRESH_TOKEN=${NEXT_PUBLIC_ENABLE_REFRESH_TOKEN}
# boilerplate:axis=auth option=jwt-refresh END
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN pnpm build

# Stage dev: local hot-reload via Tilt live_update + Turbopack HMR.
# Reuses node_modules from deps (HORIZON_KEY not needed at runtime).
# Tilt syncs src/ and public/ into /app; Turbopack detects changes and
# pushes HMR to the browser without a container restart (<1 s).
# Never used by production builds (runner is last stage).
FROM node:20-alpine AS dev
RUN corepack enable && corepack prepare pnpm@10.31.0 --activate
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NODE_ENV=development
ENV PORT=3000
ENV HOSTNAME=0.0.0.0
EXPOSE 3000
CMD ["pnpm", "dev", "--turbo"]

# Stage 3: runner
FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

COPY --from=builder --chown=nextjs:nodejs /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

CMD ["node", "server.js"]
