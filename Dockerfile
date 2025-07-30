FROM golang:1.24.5-alpine AS build

WORKDIR /usr/local/src

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN apk add --no-cache \
        curl \
    && \
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin v3.44.0 && \
    task build


FROM alpine:3

COPY --from=build /usr/local/src/dist/routing /routing

CMD [ "/routing" ]
