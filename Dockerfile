FROM busybox:ubuntu-14.04

MAINTAINER Hector Fernandez <hector@giantswarm.io>

COPY yochu /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/yochu"]
