ARG GO_VERSION=1.25.4

FROM alpine:3 AS get-task

RUN apk add --no-cache \
        curl \
    && \
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin


FROM golang:${GO_VERSION}-alpine AS build

COPY --from=get-task /usr/local/bin/task /usr/local/bin/task

WORKDIR /usr/local/src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN task build


FROM alpine:3

COPY --from=build /usr/local/src/dist/traefik-configuration-provider /usr/local/bin/traefik-configuration-provider

ENTRYPOINT [ "/usr/local/bin/traefik-configuration-provider" ]
