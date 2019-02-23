FROM ubuntu as stager

RUN mkdir /opt/goproxy && \
    echo 'geth:x:30303:' > /tmp/group && \
    echo 'geth:x:30303:30303:daemon,,,:/opt/goproxy:/bin/false' > /tmp/passwd && \
    apt-get update; apt-get install -y git golang curl jq

ENV GOPATH /opt
RUN mkdir -p /opt/src/gitlab.com/blockpane/honeywallet/goproxy
WORKDIR /opt/src/gitlab.com/blockpane/honeywallet/goproxy
COPY . /opt/src/gitlab.com/blockpane/honeywallet/goproxy
RUN go get && go build -ldflags '-s -w' -o /tmp/goproxy ./cmd/main.go

# Second stage of build, doing it this way makes for a very small container
FROM scratch

COPY --from=stager /lib64 /lib64
COPY --from=stager /lib /lib
COPY --from=stager /tmp/passwd /etc/passwd
COPY --from=stager /tmp/group /etc/group
COPY --from=stager --chown=geth /opt/goproxy /opt/goproxy
COPY --from=stager /tmp/goproxy /bin/goproxy

USER geth
WORKDIR /opt/goproxy
VOLUME /opt/goproxy/logs

ENTRYPOINT ["/bin/goproxy"]
