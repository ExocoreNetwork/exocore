FROM golang:1.21-alpine AS build
RUN apk add build-base linux-headers git
# Set working directory for the build
WORKDIR work
# Add source files
COPY . ./
# build exocore
#COPY ./exocored /go/work/build
RUN LEDGER_ENABLED=false make build

#####################################
FROM alpine AS run
RUN apk add bash curl jq
EXPOSE 26656 26657 1317 9090 8545 8546
CMD ["start", "--log_format", "plain","--chain-id","exocoretestnet_233-1", "--metrics", "--json-rpc.api", "eth,txpool,personal,net,debug,web3", "--api.enable", "--json-rpc.enable", "true"]
STOPSIGNAL SIGTERM
VOLUME /exocore
WORKDIR /exocore

COPY networks/local/exocore/wrapper.sh /usr/bin/wrapper.sh
COPY --from=build /go/work/build/exocored /exocore

ENTRYPOINT ["/usr/bin/wrapper.sh"]