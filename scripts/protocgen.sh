#!/usr/bin/env bash

# --------------
# Commands to run locally
# docker run --network host --rm -v $(CURDIR):/workspace --workdir /workspace ghcr.io/cosmos/proto-builder:v0.11.6 sh ./protocgen.sh
#
set -eo pipefail

echo "Generating gogo proto code"
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
#proto_dirs="proto/exocore/reward/v1beta1"
#echo $proto_dirs
for dir in $proto_dirs; do
	proto_files=$(find "${dir}" -maxdepth 3 -name '*.proto')
	for file in $proto_files; do
		# Check if the go_package in the file is pointing to evmos
		if grep -q "option go_package.*exocore" "$file"; then
			buf generate --template proto/buf.gen.gogo.yaml "$file"
		fi
	done
done

# move proto files to the right places
cp -r github.com/ExocoreNetwork/exocore/* ./
rm -rf github.com
