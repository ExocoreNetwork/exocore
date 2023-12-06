#!/bin/bash
current_dir="$PWD"
CHAINDIR="$current_dir/build/.testnets"

for node in {0..3}; do
    NODE_DIR="$CHAINDIR/node$node"
    GENESIS="$NODE_DIR/evmosd/config/genesis.json"
    TMP_GENESIS="$NODE_DIR/evmosd/config/tmp_genesis.json"
    APP_TOML="$NODE_DIR/evmosd/config/app.toml"
    CONFIG_TOML="$NODE_DIR/evmosd/config/config.toml"
     # If TMP_GENESIS directory does not exist, create it
        TMP_GENESIS_DIR=$(dirname "$TMP_GENESIS")
        if [ ! -d "$TMP_GENESIS_DIR" ]; then
            mkdir -p "$TMP_GENESIS_DIR"
        fi
     # used to exit on first error (any non-zero exit code)
    set -e

    # Update total supply with claim values
    # Bc is required to add this big numbers
    # total_supply=$(bc <<< "$amount_to_claim+$validators_supply")
    #total_supply=100004000000000000000010000
    #jq -r --arg total_supply "$total_supply" '.app_state.bank.supply[0].amount=$total_supply' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"


    # make sure the localhost IP is 0.0.0.0
    sed -i.bak 's/localhost/0.0.0.0/g' "$CONFIG_TOML"
    sed -i.bak 's/127.0.0.1/0.0.0.0/g' "$APP_TOML"

    # use timeout_commit 1s to make test faster
    sed -i.bak 's/timeout_commit = "3s"/timeout_commit = "1s"/g' "$CONFIG_TOML"

    # Enable the APIs for the tests to be successful
    sed -i.bak 's/enable = false/enable = true/g' "$APP_TOML"

  echo "Modified configurations for node$node"
done



