# Tiltfile — local dev orchestration for boilerplate-monorepo.
#
# Usage:
#   tilt up          # start all services with hot-reload
#   tilt down        # stop all services
#   tilt up backend-belajar-lintas-ronald  # start only the golang stack
#
# Prerequisites: Tilt >= 0.33 (https://docs.tilt.dev/install.html)
# HORIZON_KEY is auto-read from frontend/.env for the Next.js dev build; export
# HORIZON_KEY in the shell to override (see horizon_key() below).
#
# Hot-reload behaviour:
#   BE (Go):  Tilt syncs *.go files → restart_container() → CMD rebuilds binary (~2-5 s)
#   FE (Next.js): Tilt syncs src/ → Turbopack HMR → browser update (<1 s)
#
# boilerplate marker dimension: subproject
# Used by sdlc-core:scaffold-product-monorepo pruner to strip FE sections
# from BE-only projects and BE sections from FE-only projects.

# compose.dev.yml layers dev-only overrides (e.g. NODE_ENV=development for
# `next dev`) on top of the production-shaped compose.yml. The override entry is
# subproject-scoped so a BE-only prune drops it and leaves docker_compose with
# just compose.yml.
compose_files = ['compose.yml']
# boilerplate:subproject=nextjs START
compose_files.append('compose.dev.yml')
# boilerplate:subproject=nextjs END
docker_compose(compose_files)

# boilerplate:subproject=golang START
docker_build(
    'backend-belajar-lintas-ronald',
    context='./backend',
    dockerfile='backend/Dockerfile',
    target='dev',
    live_update=[
        sync('./backend/internal', '/go/src/app/internal'),
        sync('./backend/cmd',      '/go/src/app/cmd'),
        sync('./backend/pkg',      '/go/src/app/pkg'),
        sync('./backend/main.go',  '/go/src/app/main.go'),
        sync('./backend/env',      '/go/src/app/env'),
        restart_container(),
    ],
)
# boilerplate:subproject=golang END

# boilerplate:subproject=java START
docker_build(
    'template-be-java',
    context='./java',
    dockerfile='java/Dockerfile',
    live_update=[
        sync('./java/src', '/app/src'),
        restart_container(),
    ],
)
# boilerplate:subproject=java END

# boilerplate:subproject=nextjs START
# Resolve HORIZON_KEY for the Next.js dev build. Tilt's os.getenv() reads only
# the shell environment — it does NOT auto-load .env — so we read frontend/.env
# ourselves (matching how compose interpolates ${HORIZON_KEY}). Shell env wins
# when set, so `export HORIZON_KEY=...` still overrides the file.
def horizon_key():
    val = os.getenv('HORIZON_KEY', '')
    if val:
        return val
    env_path = './frontend/.env'
    if os.path.exists(env_path):
        for line in str(read_file(env_path)).splitlines():
            line = line.strip()
            if line.startswith('HORIZON_KEY='):
                return line[len('HORIZON_KEY='):].strip().strip('"').strip("'")
    return ''

docker_build(
    'frontend-belajar-lintas-ronald',
    context='./frontend',
    dockerfile='frontend/Dockerfile',
    target='dev',
    build_args={'HORIZON_KEY': horizon_key()},
    live_update=[
        sync('./frontend/src',    '/app/src'),
        sync('./frontend/public', '/app/public'),
    ],
)
# boilerplate:subproject=nextjs END
