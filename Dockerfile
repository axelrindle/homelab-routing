FROM alpine:3 AS prepare

COPY ./dist/* /usr/local/bin/traefik-configuration-provider

RUN chmod +x /usr/local/bin/traefik-configuration-provider


FROM gcr.io/distroless/static-debian13:nonroot

COPY --from=prepare /usr/local/bin/traefik-configuration-provider /usr/local/bin/traefik-configuration-provider

ENTRYPOINT [ "/usr/local/bin/traefik-configuration-provider" ]

