# Yochu

[![Build Status](https://api.travis-ci.org/giantswarm/yochu.svg)](https://travis-ci.org/giantswarm/yochu) [![](https://godoc.org/github.com/giantswarm/yochu?status.svg)](http://godoc.org/github.com/giantswarm/yochu) [![IRC Channel](https://img.shields.io/badge/irc-%23giantswarm-blue.svg)](https://kiwiirc.com/client/irc.freenode.net/#giantswarm)

`yochu` provisions already running CoreOS hosts with Docker, `etcd`, `fleet`,`rkt`, `kubectl` and `iptables`.

Host operating systems are provisioned with [giantswarm/mayu](http://github.com/giantswarm/mayu) on bare metal clusters.

## Getting Yochu

Download the latest release: https://github.com/giantswarm/yochu/releases/latest

Clone the git repository: https://github.com/giantswarm/yochu.git

## Running Yochu

Place the following unit file in your cloud-config, replacing your subnet and docker subnets:
```
[Unit]
Description=Giant Swarm Yochu
Wants=network-online.target
After=network-online.target
Before=multi-user.target
[Service]
Type=oneshot
ExecStartPre=/usr/bin/mkdir -p /home/core/bin
ExecStartPre=/usr/bin/wget --no-verbose https://downloads.giantswarm.io/yochu/0.18.0/yochu -O /home/core/bin/yochu
ExecStartPre=/usr/bin/chmod +x /home/core/bin/yochu
ExecStart=/home/core/bin/yochu setup -v -d --start-daemons=true --subnet=<your subnet> --docker-subnet=<your docker subnet> --http-endpoint=https://downloads.giantswarm.io --fleet-version=v0.11.3-gs-2 --etcd-version=v2.1.0-gs-1
RemainAfterExit=yes
[Install]
WantedBy=multi-user.target
```

## Further Steps

Check more detailed documentation: [docs](docs)

Check code documentation: [godoc](https://godoc.org/github.com/giantswarm/yochu)

Our custom binaries can be found at:
- etcd: https://downloads.giantswarm.io/etcd/v2.1.0-gs-1/etcd
- etcdctl: https://downloads.giantswarm.io/etcd/v2.1.0-gs-1/etcdctl
- fleet: https://downloads.giantswarm.io/fleet/v0.11.3-gs-2/fleetd
- fleetctl: https://downloads.giantswarm.io/fleet/v0.11.3-gs-2/fleetctl
- rkt: https://downloads.giantswarm.io/rkt/v1.1.0/rkt
- kubectl: https://downloads.giantswarm.io/k8s/v1.1.8/kubectl

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/#!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/yochu/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the
contribution workflow as well as reporting bugs.

## License

Yochu is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.

## Origin of the Name

`yochu` (ようちゅう[幼虫] pronounced "yo-choo") is Japanese for larva or chrysalis.
