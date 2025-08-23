FROM debian:stable-slim AS build

WORKDIR /usr/local/src

COPY . .

ARG DEB_VERSION
ARG DEB_ARCH

RUN echo "Architecture: ${DEB_ARCH}" >> DEBIAN/control
RUN echo "Version: ${DEB_VERSION}" >> DEBIAN/control
RUN dpkg-deb --build . "homelab-routing-${DEB_VERSION}-${DEB_ARCH}.deb"


FROM scratch

COPY --from=build /usr/local/src/*.deb .
