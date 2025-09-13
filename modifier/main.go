package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var regexApplet = regexp.MustCompile(`^(//applet:.+\(\s*APPLET[_A-Z0-9]*\()([^,]+)(, .*)$`)

/*
//usage:#define tar_trivial_usage
//usage:#define tar_full_usage "\n\n"
//usage:#define tar_example_usage
//usage:#define tree_trivial_usage NOUSAGE_STR
//usage:#define tree_full_usage ""
*/
var regexUsage = regexp.MustCompile(`^(\s*//usage:#define )([a-z0-9]+)(_[a-z0-9]+_usage\s*.*)$`)

// exit(ls_main(/*argc_unused*/ 0, (char**) argv));
// int id_main(int argc, char **argv) MAIN_EXTERNALLY_VISIBLE;
// int id_main(int argc UNUSED_PARAM, char **argv)
// int blkid_main(int argc, char **argv) MAIN_EXTERNALLY_VISIBLE;
// int blkid_main(int argc UNUSED_PARAM, char **argv)
var regexMain = regexp.MustCompile(`^(.+?)([a-z0-9]+)(_main\s*\(.+)$`)

// ret = run_applet_main(argv, kill_main);
var regexMain2 = regexp.MustCompile(`^(.+run_applet_main\(\w+, )([a-z0-9]+)(_main\);)$`)

var regexBuiltIn = regexp.MustCompile(`^(.*BLTIN\(")([a-z0-9]+)(".+)$`)

/*
alias asphalt="tar"
alias aunleitung="man"
alias ausnutza="msfconsole"
alias auspockn="unzip"
alias aussemitdi="export"
alias austausch="sed"
alias ausweis="id"
alias bam="tree"
alias beoarbeit="vim"
alias bittelassmichgehenichkommnichtraus="vim"
alias besteig="mount"
alias beweg="mv"
alias bitte="sudo"
alias boa="drill"
alias bringum="kill"
alias büdschirm="screen"
alias chef="sudo"
alias chü="sleep"
alias daumpf="steam"
alias daumpflok="sl"
alias duaher="fg"
alias duaumme="scp"
alias duaweg="bg"
alias drahdiham="rm -rf"
alias drahisgraflo="shutdown now"
alias einpockn="zip"
alias eisstockschiassn="curl"
alias erstö="mkdir"
alias fenstamoasta="tmux"
alias fernschoin="ssh"
alias finga="finger"
alias geh="sudo"
alias gehausse="cd .."
alias gehbitte="sudo"
alias geheine="cd"
alias gehgaunzausse="cd /"
alias gehzruck="cd -"
alias gehham="cd"
alias gehsterbn="kill -KILL"
alias gehscheissn="killall -9"
alias gib="curl"
alias gibher="curl"
alias greifau="touch"
alias grob="dig"
alias gruppn="groups"
alias gröwin="grep"
alias gschicht="history"
alias heast="sudo"
alias heisl="wc"
alias herdamit="apt install"
alias hobiaahnung="help"
alias hobidesnedschoamoigmocht="history"
alias hobididroht="sudo su"
alias hoiowa="wget"
alias horchzua="nc -nlvp "
alias hüfe="help"
alias ichwillnichtmehrichkannnichtmehrichhaltedasallesnichtmehraus="exit"
alias isdonofrei="df -h"
alias isschorecht="yes"
alias kloa="nano"
alias knopf="node"
alias kopf="head"
alias kopier="cp"
alias kotz="cat"
alias kua="cowsay"
alias luag="locate"
alias mehr="more"
alias moch="make"
alias mochweg="clear"
alias netz="ip"
alias netztofin="iptables"
alias netzfütatofin="nft"
alias netzkoatn="nmap"
alias netzkotz="nc"
alias noamoi="reset"
alias oanzeln="sort -u"
alias obn="top"
alias oida="sudo"
alias owaiwü="sudo"
alias passauf="watch"
alias pfiatdi="exit"
alias putz="clear"
alias qön="source"
alias ramzaum="rm"
alias klopfau="ping"
alias rechna="bc"
alias ruafau="telnet"
alias ruf="ping"
alias schaumamoi="eval"
alias schlang="python"
alias schleichdi="kill -TERM"
alias schlof="sleep"
alias schneid="cut"
alias schoin="bash"
alias schoitaus="shutdown now"
alias schrei="echo"
alias schreibauf="tcpdump"
alias schwaunz="tail"
alias sortier="sort"
alias stöei="export"
alias suach="find"
alias trogumme="scp"
alias tunnö="ssh -D 1080"
alias umbringa="kill"
alias ummödn="su"
alias undwersansiebitte="whois"
alias unterschwöllig="subl"
alias vagleich="diff"
alias varreck="kill -KILL"
alias wegdamit="apt remove"
alias weniga="less"
alias werbini="whoami"
alias werhorcht="netstat"
alias wiagehtsmeinhipsta="bashtop"
alias wiaismeibärvoamens="htop"
alias wiebitte="history"
alias wosduaido="history"
alias woisnbitte="rg -i"
alias wosisn="apropos"
alias wobini="pwd"
alias wosgeht="ps"
alias woslaft="ps"
alias wowoameileistung="htop"
alias zaumrama="rm"
alias zoag="ls"
alias zoaga="ln"
alias zoagoiss="ls -la"
alias zruck="cd -"
alias zöh="wc"
alias züntdiau="burnMMX"
*/

/*
Currently defined functions:
[, [[, acpid, add-shell, addgroup, adduser, adjtimex, arch, arp, arping, ascii, ash, asphalt, aunleitung, auspockn, austausch, ausweis, awk, bam, base32, base64, basename, bc, beep, beweg,
bittelassmichgehenichkommnichtraus, blkdiscard, blkid, blockdev, bootchartd, brctl, bunzip2, bzcat, bzip2, cal, cat, chat, chattr, chgrp, chmod, chown, chpasswd, chpst, chroot, chrt, chvt,
chü, cksum, clear, cmp, comm, conspy, cp, cpio, crc32, crond, crontab, cryptpw, cttyhack, cut, date, dc, dd, deallocvt, delgroup, deluser, depmod, devmem, df, dhcprelay, diff, dirname,
dmesg, dnsd, dnsdomainname, dos2unix, dpkg, dpkg-deb, drahdiham, du, dumpkmap, dumpleases, echo, ed, egrep, eject, env, envdir, envuidgid, erstö, ether-wake, expand, expr, factor,
fakeidentd, fallocate, false, fatattr, fbset, fbsplash, fdflush, fdformat, fdisk, fgconsole, fgrep, find, findfs, flock, fold, free, freeramdisk, fsck, fsck.minix, fsfreeze, fstrim, fsync,
ftpd, ftpget, ftpput, fuser, gehscheissn, gehsterbn, getfattr, getopt, getty, gruppn, gröwin, gunzip, gzip, halt, hd, hdparm, head, heisl, hexdump, hexedit, hostid, hostname, httpd, hush,
hwclock, i2cdetect, i2cdump, i2cget, i2cset, i2ctransfer, ifconfig, ifdown, ifenslave, ifplugd, ifup, inetd, init, insmod, install, ionice, iostat, ip, ipaddr, ipcalc, ipcrm, ipcs, iplink,
ipneigh, iproute, iprule, iptunnel, kbd_mode, killall5, klogd, last, less, link, linux32, linux64, linuxrc, ln, loadfont, loadkmap, logger, login, logname, logread, losetup, lpd, lpq, lpr,
lsattr, lsmod, lsof, lspci, lsscsi, lsusb, lzcat, lzma, lzop, makedevs, makemime, md5sum, mdev, mesg, microcom, mim, mkdosfs, mke2fs, mkfifo, mkfs.ext2, mkfs.minix, mkfs.vfat, mknod,
mkpasswd, mkswap, mktemp, modinfo, modprobe, more, mount, mountpoint, mpstat, mt, nameif, nanddump, nandwrite, nbd-client, nc, netstat, nice, nl, nmeter, nohup, nologin, nproc, nsenter,
nslookup, ntpd, od, openvt, partprobe, passwd, paste, patch, pgrep, pidof, ping, ping6, pipe_progress, pivot_root, pkill, pmap, popmaildir, poweroff, powertop, printenv, printf, ps, pscan,
pstree, pwdx, raidautorun, rdate, rdev, readahead, readlink, readprofile, realpath, reboot, reformime, remove-shell, renice, reset, resize, resume, rev, rmdir, rmmod, route, rpm, rpm2cpio,
rtcwake, run-init, run-parts, runlevel, runsv, runsvdir, rx, script, scriptreplay, seedrng, sendmail, seq, setarch, setconsole, setfattr, setfont, setkeycodes, setlogcons, setpriv, setserial,
setsid, setuidgid, sh, sha1sum, sha256sum, sha384sum, sha3sum, sha512sum, showkey, shred, shuf, slattach, smemcap, softlimit, sort, split, ssl_client, start-stop-daemon, stat, strings, stty,
su, sulogin, sum, sv, svc, svlogd, svok, swapoff, swapon, switch_root, sync, sysctl, syslogd, tac, tail, taskset, tcpsvd, tee, telnet, telnetd, test, tftp, tftpd, time, timeout, top, touch,
tr, traceroute, traceroute6, true, truncate, ts, tsort, tty, ttysize, tunctl, ubiattach, ubidetach, ubimkvol, ubirename, ubirmvol, ubirsvol, ubiupdatevol, udhcpc, udhcpc6, udhcpd, udpsvd,
uevent, umount, uname, unexpand, uniq, unix2dos, unlink, unlzma, unshare, unxz, uptime, users, usleep, uudecode, uuencode, vconfig, vlock, volname, w, wall, watch, watchdog, wget, which, who,
whoami, whois, wobini, xargs, xxd, xz, xzcat, yes, zcat, zcip, zoag
*/

var mapping = map[string]string{
	"ls":       "zoag",
	"tar":      "asphalt",
	"man":      "aunleitung",
	"unzip":    "auspockn",
	"sed":      "austausch",
	"id":       "ausweis",
	"groups":   "gruppn",
	"tree":     "bam",
	"vi":       "bittelassmichgehenichkommnichtraus",
	"mv":       "beweg",
	"sudo":     "chef",
	"sleep":    "chü",
	"steam":    "daumpf",
	"sl":       "daumpflok",
	"fg":       "duaher",
	"scp":      "duaumme",
	"bg":       "duaweg",
	"rm":       "drahdiham",
	"shutdown": "drahisgraflo",
	"zip":      "einpockn",
	"curl":     "eisstockschiassn",
	"mkdir":    "erstö",
	"tmux":     "fenstamoasta",
	"ssh":      "fernschoin",
	"finger":   "finga",
	"cd":       "geheine",
	"kill":     "gehsterbn",
	"killall":  "gehscheissn",
	"dig":      "grob",
	"grep":     "gröwin",
	"history":  "gschicht",
	"wc":       "heisl",
	"pwd":      "wobini",
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	pathIn := flag.String("path", "", "path to the directory to scan")
	flag.Parse()
	if *pathIn == "" {
		return errors.New("please provide a valid path using the -path flag")
	}

	err := filepath.Walk(*pathIn, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		ext := filepath.Ext(path)
		if !info.IsDir() && (ext == ".c" || ext == ".h") {
			if err := processFile(path); err != nil {
				return fmt.Errorf("error processing file %q: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path %q: %w", ".", err)
	}

	return nil
}

func processFile(path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var newContents []byte

	scanner := bufio.NewScanner(bytes.NewReader(contents))
	for scanner.Scan() {
		t := scanner.Text()
		switch {
		case regexApplet.MatchString(t):
			modified := process(path, t, regexApplet)
			t = modified
		case regexUsage.MatchString(t):
			modified := process(path, t, regexUsage)
			t = modified
		case regexMain2.MatchString(t):
			modified := process(path, t, regexMain2)
			t = modified
		case regexMain.MatchString(t):
			modified := process(path, t, regexMain)
			t = modified
		case regexBuiltIn.MatchString(t):
			modified := process(path, t, regexBuiltIn)
			t = modified
		}
		newContents = append(newContents, []byte(t+"\n")...)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error occurred: %w", err)
	}

	if err := os.WriteFile(path, newContents, 0644); err != nil { // nolint:gosec
		return err
	}

	return nil
}

func process(path string, line string, r *regexp.Regexp) string {
	return r.ReplaceAllStringFunc(line, func(match string) string {
		matches := r.FindStringSubmatch(match)
		if len(matches) == 4 {
			prefix := matches[1]   // Everything before the name
			original := matches[2] // The original name (first capture group)
			suffix := matches[3]   // Everything after the name

			// handle special case where another command is a variation of the main command
			// and references it in the parameters. As we only process one line per default,
			// these need to be handled here.
			switch original {
			case "groups":
				//applet:IF_GROUPS(APPLET_NOEXEC(groups, id, BB_DIR_USR_BIN, BB_SUID_DROP, groups))
				//applet:IF_ID(    APPLET_NOEXEC(id,     id, BB_DIR_USR_BIN, BB_SUID_DROP, id    ))
				newKeyword, ok := mapping["id"]
				if ok {
					suffix = strings.ReplaceAll(suffix, ", id,", fmt.Sprintf(", %s,", newKeyword))
				}
			case "killall", "killall5":
				//applet:IF_KILL(    APPLET_NOFORK(kill,     kill, BB_DIR_BIN,      BB_SUID_DROP, kill))
				//                   APPLET_NOFORK:name      main  location         suid_type     help
				//applet:IF_KILLALL( APPLET_NOFORK(killall,  kill, BB_DIR_USR_BIN,  BB_SUID_DROP, killall))
				//applet:IF_KILLALL5(APPLET_NOFORK(killall5, kill, BB_DIR_USR_SBIN, BB_SUID_DROP, killall5))
				newKeyword, ok := mapping["kill"]
				if ok {
					suffix = strings.ReplaceAll(suffix, ", kill,", fmt.Sprintf(", %s,", newKeyword))
					suffix = strings.ReplaceAll(suffix, ",  kill,", fmt.Sprintf(", %s,", newKeyword))
				}
			case "egrep", "fgrep":
				//applet:IF_EGREP(APPLET_ODDNAME(egrep, grep, BB_DIR_BIN, BB_SUID_DROP, egrep))
				//applet:IF_FGREP(APPLET_ODDNAME(fgrep, grep, BB_DIR_BIN, BB_SUID_DROP, fgrep))
				newKeyword, ok := mapping["grep"]
				if ok {
					suffix = strings.ReplaceAll(suffix, ", grep,", fmt.Sprintf(", %s,", newKeyword))
				}
			}

			if newName, exists := mapping[original]; exists {
				// also replace parameter names in suffix
				var r2 = regexp.MustCompile(fmt.Sprintf(`,\s*%s`, original))
				suffix = r2.ReplaceAllString(suffix, fmt.Sprintf(", %s", newName))

				modifiedLine := prefix + newName + suffix
				fmt.Printf("%s: %s\n", path, modifiedLine) // nolint:forbidigo
				return modifiedLine
			}
			fmt.Printf("%s: %s\n", path, prefix+original+suffix) // nolint:forbidigo
			return prefix + original + suffix                    // return maybe modified suffix
		}
		return match
	})
}
