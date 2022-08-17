#!/bin/bash

source get_ip.sh
SERVICE_NAME=go-micro-tpl-http
export SERVICE_ID_HTTP=${SERVICE_NAME}-${IP_STR}
export SERVICE_HTTP_CONFIG=`cat <<EOF
{
  "ID": "${SERVICE_ID_HTTP}",
  "name": "${SERVICE_NAME}",
  "tags": [],
  "port": ${HTTP_PORT},
  "check": {
    "http": "http://localhost:${HTTP_PORT}/healthz",
    "interval": "10s",
    "timeout": "1s"
  }
}
EOF`