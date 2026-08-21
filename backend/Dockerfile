# Stage build: will build the go into single binary
FROM golang:1.25-bookworm AS build
WORKDIR /go/src/app
COPY ./cmd ./cmd
COPY ./db ./db
COPY ./env ./env
COPY ./internal ./internal
# boilerplate:axis=rpc option=grpc START
# Committed gRPC stubs imported by internal/app/container.go. Required for the
# build stage; the dev/final stages inherit the compiled binary so need no copy.
COPY ./gen ./gen
# boilerplate:axis=rpc option=grpc END
COPY ./pkg ./pkg
COPY ./main.go ./
COPY ./go.mod ./
COPY ./go.sum ./

RUN go env -w CGO_ENABLED=0
RUN go env -w GO111MODULE="on"
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

# Install dependencies using docker cache.
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build binary using docker cache.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o backend .

# Stage dev: local hot-reload via Tilt live_update + restart_container().
# Tilt syncs source files into /go/src/app, then calls restart_container()
# which re-runs this CMD: Go rebuilds the binary incrementally (~2-5 s)
# and starts the server. Never used by production builds (final is last stage).
#
# File layout mirrors the final stage (/backend, /env, /db at root) so that
# backend-belajar-lintas-ronald-migrate and backend-belajar-lintas-ronald-seed (which share this image and
# run ["/backend", "db:migrate"] / ["/backend", "db:seed"]) work correctly.
FROM build AS dev
COPY --from=build /go/src/app/backend /backend
COPY --from=build /go/src/app/env /env
COPY --from=build /go/src/app/db /db
# Download modules into the image layer so CMD go build never needs network access.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
CMD ["sh", "-c", "go build -o /backend . && /backend serve"]

# Stage final: the final runner of the image
FROM gcr.io/distroless/static-debian10 AS final
USER nonroot

COPY --from=build /go/src/app/backend /
COPY --from=build /go/src/app/env /env
COPY --from=build /go/src/app/db /db

CMD ["/backend", "serve"]
