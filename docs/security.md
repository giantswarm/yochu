# Security Overview

## HTTPS Endpoints

Yochu uses TLS to fetch custom binaries either from S3 or from a running instance of [Mayu](https://github.com/giantswarm/mayu).

However, TLS does not ensure that the binaries have not been tinkered with. Thus, if in doubt, we recommend check binaries before and after download to ensure they have not been modified.

## iptables

If enabled, Yochu deploys a custom set of `iptables` rules to the host. Following rules are deployed:
* Allow NATing and masquerade traffic from the docker subnet that isn't sent to another container on the host.
* Default policy is DROP
* Allow destination subnet, if the packets are not coming from the docker bridge
* Allow established connections
* Forward established connections to containers
* Forward packets from containers to default gateway
* Don't forward packets from containers to other machines in the private subnet
* Forward packets from containers that aren't sent over docker0 (e.g. for Weave Net)
* Forward packets over docker0 to DOCKER chain (e.g. all incoming packets for docker portmapping)
* Allow all outgoing packets that haven't been dropped by rules before

As Yochu can set up basic configuration of e.g. network on hosts, we recommend to thoroughly review the configurations you are deploying before enabling this option.
