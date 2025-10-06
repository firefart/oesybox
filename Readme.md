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

## Commmands supported

```text
BusyBox v1.38.0.git (2025-10-06 00:35:44 UTC) multi-call binary.
BusyBox is copyrighted by many authors between 1998-2015.
Licensed under GPLv2. See source distribution for detailed
copyright notices.

Usage: busybox [function [arguments]...]
   or: busybox --list[-full]
   or: busybox --show SCRIPT
   or: busybox --install [-s] [DIR]
   or: function [arguments]...

	BusyBox is a multi-call binary that combines many common Unix
	utilities into a single executable.  Most people will create a
	link to busybox for each function they wish to use and BusyBox
	will act like whatever it was invoked as.

Currently defined functions:
	[, [[, acpid, add-shell, addgroup, adduser, adjtimex, arch, arp,
	arping, ascii, ash, asphalt, aunleitung, auspockn, austausch, ausweis,
	awk, bam, base32, base64, basename, beep, besteig, beweg,
	bittelassmichgehenichkommnichtraus, blkdiscard, blkid, blockdev,
	bootchartd, brctl, bunzip2, bzcat, bzip2, cal, chat, chattr, chgrp,
	chmod, chown, chpasswd, chpst, chroot, chrt, chvt, chü, cksum, cmp,
	comm, conspy, cpio, crc32, crond, crontab, cryptpw, cttyhack, date, dc,
	dd, deallocvt, delgroup, deluser, depmod, devmem, df, dhcprelay,
	dirname, dmesg, dnsd, dnsdomainname, dos2unix, dpkg, dpkg-deb,
	drahdiham, du, dumpkmap, dumpleases, ed, egrep, eject, env, envdir,
	envuidgid, erstö, ether-wake, expand, expr, factor, fakeidentd,
	fallocate, false, fatattr, fbset, fbsplash, fdflush, fdformat, fdisk,
	fgconsole, fgrep, findfs, flock, fold, free, freeramdisk, fsck,
	fsck.minix, fsfreeze, fstrim, fsync, ftpd, ftpget, ftpput, fuser,
	gehscheissn, gehsterbn, getfattr, getopt, getty, greifau, gruppn,
	gröwin, gunzip, gzip, halt, hd, hdparm, heisl, hexdump, hexedit,
	hoiowa, hostid, hostname, httpd, hush, hwclock, i2cdetect, i2cdump,
	i2cget, i2cset, i2ctransfer, ifconfig, ifdown, ifenslave, ifplugd,
	ifup, inetd, init, insmod, install, ionice, iostat, ip, ipaddr, ipcalc,
	ipcrm, ipcs, iplink, ipneigh, iproute, iprule, iptunnel, kbd_mode,
	killall5, klogd, klopfau, kopf, kopier, kotz, last, link, linux32,
	linux64, linuxrc, loadfont, loadkmap, logger, login, logname, logread,
	losetup, lpd, lpq, lpr, lsattr, lsmod, lsof, lspci, lsscsi, lsusb,
	lzcat, lzma, lzop, makedevs, makemime, md5sum, mdev, mehr, mesg,
	microcom, mim, mkdosfs, mke2fs, mkfifo, mkfs.ext2, mkfs.minix,
	mkfs.vfat, mknod, mkpasswd, mkswap, mktemp, mochweg, modinfo, modprobe,
	mountpoint, mpstat, mt, nameif, nanddump, nandwrite, nbd-client,
	netzkotz, nice, nl, nmeter, noamoi, nohup, nologin, nproc, nsenter,
	nslookup, ntpd, od, openvt, partprobe, passauf, passwd, paste, patch,
	pgrep, pidof, ping6, pipe_progress, pivot_root, pkill, pmap,
	popmaildir, poweroff, powertop, printenv, printf, ps, pscan, pstree,
	pwdx, raidautorun, rdate, rdev, readahead, readlink, readprofile,
	realpath, reboot, rechna, reformime, remove-shell, renice, resize,
	resume, rev, rmdir, rmmod, route, rpm, rpm2cpio, rtcwake, ruafau,
	run-init, run-parts, runlevel, runsv, runsvdir, rx, schneid, schrei,
	schwaunz, script, scriptreplay, seedrng, sendmail, seq, setarch,
	setconsole, setfattr, setfont, setkeycodes, setlogcons, setpriv,
	setserial, setsid, setuidgid, sh, sha1sum, sha256sum, sha384sum,
	sha3sum, sha512sum, showkey, shred, shuf, slattach, smemcap, softlimit,
	sortier, split, ssl_client, start-stop-daemon, stat, strings, stty, su,
	suach, sulogin, sum, sv, svc, svlogd, svok, swapoff, swapon,
	switch_root, sync, sysctl, syslogd, tac, taskset, tcpsvd, tee, telnetd,
	test, tftp, tftpd, time, timeout, top, tr, traceroute, traceroute6,
	true, truncate, ts, tsort, tty, ttysize, tunctl, ubiattach, ubidetach,
	ubimkvol, ubirename, ubirmvol, ubirsvol, ubiupdatevol, udhcpc, udhcpc6,
	udhcpd, udpsvd, uevent, umount, uname, undwersansiebitte, unexpand,
	uniq, unix2dos, unlink, unlzma, unshare, unxz, uptime, users, usleep,
	uudecode, uuencode, vagleich, vconfig, vlock, volname, w, wall,
	watchdog, weniga, werbini, werhorcht, which, who, wobini, xargs, xxd,
	xz, xzcat, yes, zcat, zcip, zoag, zoaga

```