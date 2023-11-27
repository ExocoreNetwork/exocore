#!/usr/bin/env sh
set -euo pipefail
set -x

BINARY=/exocore/${BINARY:-evmosd}
ID=${ID:-0}
LOG=${LOG:-exocore.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'evmosd'"
	exit 1
fi

export EHOME="/data/node${ID}/evmosd"

if [ -d "$(dirname "${EHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${EHOME}" "$@" | tee "${EHOME}/${LOG}"
else
  "${BINARY}" --home "${EHOME}" "$@"
fi
