#!/bin/sh

#this needs to be a certain size for higher bandwidths to work properly
export TRAFFIC_CONTROL_BURST_BUFFER="16mbit"

if [ $ENABLE_TRAFFIC_CONTROL == 'yes' ]
then
	tc qdisc add dev eth0 parent root handle 1: hfsc default 11
	tc class add dev eth0 parent 1: classid 1:1 hfsc sc rate ${OUTGOING_BANDWIDTH} ul rate ${OUTGOING_BANDWIDTH}
	tc class add dev eth0 parent 1:1 classid 1:11 hfsc sc rate ${OUTGOING_BANDWIDTH}
	tc qdisc add dev eth0 handle ffff: ingress
	tc filter add dev eth0 parent ffff: protocol ip prio 1 u32 match ip src 0.0.0.0/0 police rate ${INCOMING_BANDWIDTH} burst ${TRAFFIC_CONTROL_BURST_BUFFER} drop flowid :1
fi

./filebrowser