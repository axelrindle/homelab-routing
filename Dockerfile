FROM gcr.io/distroless/static-debian13:nonroot

COPY ./dist/* /usr/local/bin/traefik-configuration-provider

RUN chmod +x /usr/local/bin/traefik-configuration-provider

ENTRYPOINT [ "/usr/local/bin/traefik-configuration-provider" ]
