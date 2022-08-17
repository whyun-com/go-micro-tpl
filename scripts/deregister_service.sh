#!/bin/bash
# set -x

source service_config_grpc.sh
source service_config_http.sh
deregisterService() {
    curl --request PUT http://$CONSUL_ADDR/v1/agent/service/deregister/$1
}

if [ "$GRPC_PORT" != "" ] ; then
    deregisterService "$SERVICE_ID_GRPC"
    echo deregiser grpc service finished
else
    echo GRPC_PORT not defined
fi


if [ "$HTTP_PORT" != "" ] ; then
    deregisterService "$SERVICE_ID_HTTP"
    echo deregiser http service finished
else
    echo HTTP_PORT not defined
fi