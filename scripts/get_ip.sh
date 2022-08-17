#!/bin/bash

if [ "$USER_DEFINED_IP" != "" ] ; then
  export IP=$USER_DEFINED_IP
else
  export IP=`ifconfig eth0 | grep "inet " | awk '{print $2}'`
fi

export IP_STR=${IP//\./_}