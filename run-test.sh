#!/bin/bash

PAYLOAD_DIR="test-payload-data"
N_MINER="2"

usage() {
  echo "Usage: $0 [-m <number of miners>]" 1>&2
  exit 1
}

while getopts ":m:" o; do
    case "${o}" in
        m)
            N_MINER=${OPTARG}
            [[ "$N_MINER" =~ ^[0-9]+$ ]] || usage
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

for _ in $(seq 1 "${N_MINER}"); do
  go run miner.go &
done

for payload in ${PAYLOAD_DIR}/*; do
	go run client.go -payload "$(cat "${payload}")" &
done

go run server.go -target 00000
