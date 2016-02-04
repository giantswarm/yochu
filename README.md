# Yochu

`yochu` provisions already running hosts with docker, etcd, fleet, iptables ... settings.

Hosts are launched with [giantswarm/primer](http://giantswarm/primer) for an AWS cluster.

Hosts operating system are provisioned with [giantswarm/mayu](http://github.com/giantswarm/mayu) for a bare metal clusters.

## Releasing

```
builder release minor -p
git checkout <new tag>
make publish
```

## Custom etcd and fleet binaries:
Our custom binaries can be found at:
etcd: https://downloads.giantswarm.io/etcd/v2.1.0-gs-1/etcd, https://downloads.giantswarm.io/etcd/v2.1.0-gs-1/etcdctl
fleet: https://downloads.giantswarm.io/fleet/v0.11.3-gs-2/fleetd, https://downloads.giantswarm.io/fleet/v0.11.3-gs-2/fleetctl