# backport binary
Backporting binaries enables you to pin certain binaries during the
provisioning process. This means you can setup your CoreOS cluster with
specific versions of the `docker`, `fleet`, or `etcd` binary. This is how it
looks like to pin fleet to the version currently provided by GiantSwarm. Simply
set the following flags to your `yochu` executable.

```
--http-endpoint=https://downloads.giantswarm.io --fleet-version=v0.11.5-gs-grpc-1 --steps=distribution,overlay,fleet
```

The `--http-endpoint` defaults to `https://downloads.giantswarm.io`. We are
using it here for completeness. This can be any other location though. The
expected locations for `fleet` version `v0.11.5-gs-grpc-1` would be the
following.

- `https://downloads.giantswarm.io/fleet/v0.11.5-gs-grpc-1/fleetd`
- `https://downloads.giantswarm.io/fleet/v0.11.5-gs-grpc-1/fleetctl`

See:
 * https://github.com/giantswarm/yochu/blob/master/cli/setup_cmd.go#L36
 * https://github.com/giantswarm/yochu/blob/master/steps/fleet/fleet.go#L34

Since naming is hard here a short explanation of the steps.
- Step `distribution` step prepares the readonly file system on CoreOS by
  creating certain directories.
- Step `overlay` creates the overlay mount at `/usr/bin`. So in the end your
  custom `fleetd` will be available through `/usr/bin/fleetd`.
- Step `fleet` downloads the fleet binary provided by `--fleet-version`. To
  make this step working, you need to also provide steps `distribution` and
  `overlay`. This fetches the fleet binary to the distribution path which
  defaults to `/opt/giantswarm/bin`.
