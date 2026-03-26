FROM golang:latest AS go-builder
WORKDIR /go/src/app
COPY modifier/go.* modifier/*.go /go/src/app/
ENV CGO_ENABLED=0
RUN go build -a -o modifier -ldflags="-s -w" -trimpath

FROM ubuntu:latest AS builder
RUN apt-get update && apt-get install -y --no-install-recommends \
  git \
  curl \
  ca-certificates \
  build-essential
# busybox git is down all the time and the autobuilds fail
# RUN git clone --depth 1 https://git.busybox.net/busybox/ /opt/busybox
RUN curl -L https://busybox.net/downloads/busybox-snapshot.tar.bz2 -o /tmp/busybox.tar.bz2 \
  && tar -xjf /tmp/busybox.tar.bz2 -C /opt \
  && rm -f /tmp/busybox.tar.bz2
WORKDIR /opt/busybox
COPY --from=go-builder /go/src/app/modifier /opt/modifier
RUN /opt/modifier -path /opt/busybox
# Set   # CONFIG_TC is not defined
# https://lists.busybox.net/pipermail/busybox-cvs/2024-January/041752.html
RUN make defconfig && \
  sed -i 's/CONFIG_TC=y/# CONFIG_TC is not set/' .config && \
  sed -i 's/# CONFIG_STATIC is not set/CONFIG_STATIC=y/' .config && \
  #cat .config && \
  make -j"$(nproc)"

FROM ubuntu:latest AS rootfs
RUN mkdir /rootfs && \
  mkdir /rootfs/bin && \
  mkdir /rootfs/etc && \
  mkdir /rootfs/root && \
  echo "root:x:0:0:root:/root:/bin/sh" > /rootfs/etc/passwd && \
  echo "root:x:0:" > /rootfs/etc/group

FROM scratch
# FROM debian:stable-slim
COPY --from=rootfs /rootfs/ /
COPY --from=builder /opt/busybox/busybox /bin
RUN ["/bin/busybox", "--install", "-s", "/bin"]
ENTRYPOINT ["/bin/busybox", "ash"]
