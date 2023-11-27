FROM golang:1.20-alpine AS build

RUN apk add build-base  linux-headers


# Set working directory for the build
WORKDIR work

# Add source files
COPY .  ./
# build Ethermint
RUN LEDGER_ENABLED=false make  build

FROM alpine AS run
RUN apk add bash curl jq
EXPOSE 26656 26657 1317 9090 8545 8546
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start", "--log_format", "plain"]
STOPSIGNAL SIGTERM
VOLUME /exocore
WORKDIR /exocore

COPY networks/local/exocore/wrapper.sh /usr/bin/wrapper.sh
COPY --from=build /work/build/evmosd /exocore/
