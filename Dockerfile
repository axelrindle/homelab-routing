FROM alpine:3 AS prepare

COPY ./dist/* /usr/local/bin/traefik-forwarder

RUN chmod +x /usr/local/bin/traefik-forwarder


FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=prepare /usr/local/bin/traefik-forwarder /usr/local/bin/traefik-forwarder

ENTRYPOINT [ "/usr/local/bin/traefik-forwarder" ]
