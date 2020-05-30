# Build Stage
FROM mohonishc/moonshot-backend-starter:1.13 AS build-stage

LABEL app="build-moonshot-backend"
LABEL REPO="https://github.com/mohonish/moonshot-backend"

ENV PROJPATH=/go/src/github.com/mohonish/moonshot-backend

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/mohonish/moonshot-backend
WORKDIR /go/src/github.com/mohonish/moonshot-backend

RUN make build-alpine

# Final Stage
FROM mohonishc/moonshot-backend-starter

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/mohonish/moonshot-backend"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/moonshot-backend/bin

WORKDIR /opt/moonshot-backend/bin

COPY --from=build-stage /go/src/github.com/mohonish/moonshot-backend/bin/moonshot-backend /opt/moonshot-backend/bin/
RUN chmod +x /opt/moonshot-backend/bin/moonshot-backend

# Create appuser
RUN adduser -D -g '' moonshot-backend
USER moonshot-backend

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/moonshot-backend/bin/moonshot-backend"]
