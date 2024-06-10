FROM golang:1.21.11-alpine3.18 AS build
RUN apk add --no-cache build-base=0.5-r3 linux-headers=6.3-r0 git=2.40.1-r0
# Set working directory for the build
WORKDIR /go/work
# Add source files
COPY . ./
# build exocore
#COPY ./exocored /go/work/build
RUN LEDGER_ENABLED=false make build

#####################################
FROM alpine:3.18 AS run
RUN apk add --no-cache libstdc++=12.2.1_git20220924-r10 bash=5.2.15-r5 curl=8.5.0-r0 jq=1.6-r4 \
    && addgroup -g 1000 exocore \
    && adduser -S -h /home/exocore -D exocore -u 1000 -G exocore
EXPOSE 26656 26657 1317 9090 8545 8546
# TODO: exocore testnet chainid is still under consideration and need to be finalized later
CMD ["start", "--log_format", "plain", "--chain-id", "exocoretestnet_233-1", "--metrics", "--json-rpc.api", "eth,txpool,personal,net,debug,web3", "--api.enable", "--json-rpc.enable", "true", "--minimum-gas-prices", "0.0001aexocore"]
# by default, a SIGKILL is sent after 10 seconds. We need to override this to allow graceful shutdown.
STOPSIGNAL SIGTERM
VOLUME /exocore
WORKDIR /exocore

COPY ./networks/local/exocore/wrapper.sh /usr/bin/wrapper.sh
COPY --from=build /go/work/build/exocored /exocore

# Use the created non-root user
USER exocore

HEALTHCHECK --interval=30s --timeout=30s --retries=3 CMD curl -f http://localhost:26657/health || exit 1

ENTRYPOINT ["/usr/bin/wrapper.sh"]