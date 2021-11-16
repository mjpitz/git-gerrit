FROM alpine:3.12

WORKDIR /opt/git-gerrit
RUN apk add -U tzdata ca-certificates

ARG TARGETOS
ARG TARGETARCH
COPY dist/git-gerrit_${TARGETOS}_${TARGETARCH}/git-gerrit /opt/git-gerrit/bin/git-gerrit
COPY LICENSE /opt/git-gerrit/LICENSE

ENTRYPOINT [ "/opt/git-gerrit/bin/git-gerrit" ]
