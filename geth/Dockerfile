FROM ubuntu as stager

RUN mkdir /opt/geth && \
    echo 'geth:x:30303:' > /tmp/group && \
    echo 'geth:x:30303:30303:daemon,,,:/opt/geth:/bin/false' > /tmp/passwd && \
    apt-get update; apt-get install -y git golang curl jq

ENV GOPATH /opt
RUN mkdir -p /opt/src/github.com/ethereum
WORKDIR /opt/src/github.com/ethereum
RUN curl -sL $(curl -s https://api.github.com/repos/ethereum/go-ethereum/releases/latest |jq -r .tarball_url) >go-ethereum.tgz && \
    tar zxvf go-ethereum.tgz && mv ethereum-go-ethereum-* go-ethereum
WORKDIR /opt/src/github.com/ethereum/go-ethereum
RUN make geth

# Second stage of build, doing it this way makes for a very small container
FROM scratch

COPY --from=stager /lib64 /lib64
COPY --from=stager /lib /lib
COPY --from=stager /tmp/passwd /etc/passwd
COPY --from=stager /tmp/group /etc/group
COPY --from=stager --chown=geth /opt/geth /opt/geth
COPY --from=stager /opt/src/github.com/ethereum/go-ethereum/build/bin/geth /bin/geth

USER geth
WORKDIR /opt/geth

ENTRYPOINT ["/bin/geth","--datadir","/opt/geth", \
            "--rpc","--rpcaddr","0.0.0.0","--rpcapi","personal,eth,web3", \
            "--rpcvhosts","*","--maxpeers","100","--syncmode"]
