#!/bin/bash

source get_ip.sh
SERVICE_NAME=go-micro-tpl-grpc
# set -x
export SERVICE_ID_GRPC=${SERVICE_NAME}-${IP_STR}
export SERVICE_GRPC_CONFIG=`cat <<EOF
{
  "ID": "${SERVICE_ID_GRPC}",
  "name": "${SERVICE_NAME}",
  "tags": [],
  "port": ${GRPC_PORT},
  "check": {
    "grpc": "localhost:${GRPC_PORT}/${SERVICE_NAME}",
    "interval": "10s",
    "timeout": "1s"
  }
}
EOF`