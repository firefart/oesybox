FROM golang:latest AS go-builder
WORKDIR /go/src/app
COPY modifier/go.* modifier/*.go /go/src/app/
ENV CGO_ENABLED=0
RUN go build -a -o modifier -ldflags="-s -w" -trimpath

FROM ubuntu:latest AS builder
RUN apt-get update && apt-get install -y --no-install-recommends \
  git \
  build-essential
RUN git clone --depth 1 git://busybox.net/busybox.git /opt/busybox
WORKDIR /opt/busybox
COPY --from=go-builder /go/src/app/modifier /opt/modifier
RUN /opt/modifier -path /opt/busybox
# Set   # CONFIG_TC is not defined
# https://lists.busybox.net/pipermail/busybox-cvs/2024-January/041752.html
RUN make defconfig && \
  sed -i 's/CONFIG_TC=y/# CONFIG_TC is not set/' .config && \
  sed -i 's/# CONFIG_STATIC is not set/CONFIG_STATIC=y/' .config && \
  #cat .config && \
  make -j$(nproc)

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
