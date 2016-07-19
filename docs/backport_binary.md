# backport binary
Backporting binaries enables you to pin certain binaries during the
provisioning process. This means you can setup your CoreOS cluster with
specific versions of the `docker`, `fleet`, or `etcd` binary. This is how it
looks like to pin fleet to the version currently provided by GiantSwarm. Simply
set the following flags to your `yochu` executable.

```
--fleet-version=v0.11.5-gs-grpc-1 --steps=distribution,overlay,fleet
```

Since naming is hard here a short explanation of the steps.
- Step `distribution` step prepares the readonly file system on CoreOS by
  creating certain directories.
- Step `overlay` creates the overlay file system and links folders in which
  binaries are placed in.
- Step `fleet` downloads the fleet binary provided by `--fleet-version`. To
  make this step working, you need to also provide steps `distribution` and
  `overlay`.
