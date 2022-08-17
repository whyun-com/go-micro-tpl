#!/bin/bash
# set -x

source service_config_grpc.sh
source service_config_http.sh
registerService() {
    curl http://$CONSUL_ADDR/v1/agent/service/register \
    --request PUT \
    --data "$1"
}

if [ "$GRPC_PORT" != "" ] ; then
    registerService "$SERVICE_GRPC_CONFIG"
    echo regiser grpc service finished
else
    echo GRPC_PORT not defined
fi


if [ "$HTTP_PORT" != "" ] ; then
    registerService "$SERVICE_HTTP_CONFIG"
    echo regiser http service finished
else
    echo HTTP_PORT not defined
fi