# Frontend build
FROM node:18 as frontendbuild

WORKDIR /app/frontend

COPY . /app

RUN npm ci

RUN npm run build


# Backend build
FROM golang:1.20 as build

WORKDIR /app

COPY --from=frontendbuild /app /app

RUN go mod download

RUN CGO_ENABLED=0 go build


# Runtime
FROM alpine:latest
RUN apk --update add ca-certificates \
                     mailcap \
                     curl \
                     jq \
                     iproute2-tc

COPY --from=build /app/healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

HEALTHCHECK --start-period=2s --interval=5s --timeout=3s \
    CMD /healthcheck.sh || exit 1

VOLUME /srv
EXPOSE 80

ENV ENABLE_TRAFFIC_CONTROL="no"

COPY --from=build /app/scripts/init.sh /init.sh
RUN chmod +x /init.sh
COPY --from=build /app/docker_config.json /.filebrowser.json
COPY --from=build /app/filebrowser /filebrowser

ENTRYPOINT [ "/init.sh" ]
