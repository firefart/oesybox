# OESYBOX

[oeshell](https://github.com/martinhaunschmid/oeshell/) with busybox inside docker

## Running

```bash
docker run --rm -it --pull always firefart/oesybox
```

```bash
docker run --rm -it --pull always ghcr.io/firefart/oesybox:latest
```

## Test locally

```bash
docker build --pull -t test -f Dockerfile .
docker run --rm -it test
```
