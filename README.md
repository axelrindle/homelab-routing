# homelab-gateway

> Dynamically expose Traefik routers to upstream Traefik instances

> [!NOTE]
> This project is intended for use with the [Kubernetes Ingress Controller](https://doc.traefik.io/traefik/providers/kubernetes-ingress/) (default in [K3S](https://k3s.io/)).

## About

In my homelab there are two networks:

- main network at `192.168.178.0/24`
- k3s network at `10.10.10.0/24`

Within the `main` network is a Synology NAS and a bare-metal server.
Within the `k3s` network is a K3S cluster with several other services.

All of those services are available within the homelab network using static routing and a `.local` domain.

Selectively exposing some services on an external domain (e.g. `example.dev`) turned out to be
non-trivial because of the different networks within the homelab.

Solving the problem involves a Traefik instance within the `main` network which is exposed to the internt. Services in the `main` network can be easily exposed by providing a dynamic file config
to Traefik.

This project provides a dynamic configuration endpoint for Traefik selectively
exposing routers from the `k3s` internal Traefik instance.

## Configuration

Have a look in the `example/` and `pkg/debian/` directories for example configuration.
