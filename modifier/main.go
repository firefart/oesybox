package main

import (
	"bufio"
	"bytes"
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
// ret = run_applet_main(argv, kill_main);
var regexMain = regexp.MustCompile(`^(.+?)([a-z0-9]+)(_main\s*\(.+)$`)
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
	pathIn := flag.String("path", "", "path to the directory to scan")
	flag.Parse()
	if *pathIn == "" {
		fmt.Println("Please provide a valid path using the -path flag.")
		return
	}

	err := filepath.Walk(*pathIn, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		ext := filepath.Ext(path)
		if !info.IsDir() && (ext == ".c" || ext == ".h") {
			processFile(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", ".", err)
		return
	}
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
		fmt.Printf("error occurred: %v\n", err)
	}

	if err := os.WriteFile(path, newContents, 0644); err != nil {
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

			// handle special case of
			//applet:IF_GROUPS(APPLET_NOEXEC(groups, id, BB_DIR_USR_BIN, BB_SUID_DROP, groups))
			//applet:IF_ID(    APPLET_NOEXEC(id,     id, BB_DIR_USR_BIN, BB_SUID_DROP, id    ))
			// where groups reference id
			if original == "groups" {
				new, ok := mapping["id"]
				if ok {
					suffix = strings.ReplaceAll(suffix, ", id,", fmt.Sprintf(", %s,", new))
				}
			}

			// handle special case
			//applet:IF_KILL(    APPLET_NOFORK(kill,     kill, BB_DIR_BIN,      BB_SUID_DROP, kill))
			//                   APPLET_NOFORK:name      main  location         suid_type     help
			//applet:IF_KILLALL( APPLET_NOFORK(killall,  kill, BB_DIR_USR_BIN,  BB_SUID_DROP, killall))
			//applet:IF_KILLALL5(APPLET_NOFORK(killall5, kill, BB_DIR_USR_SBIN, BB_SUID_DROP, killall5))
			if original == "killall" || original == "killall5" {
				new, ok := mapping["kill"]
				if ok {
					suffix = strings.ReplaceAll(suffix, ", kill,", fmt.Sprintf(", %s,", new))
					suffix = strings.ReplaceAll(suffix, ",  kill,", fmt.Sprintf(", %s,", new))
				}
			}

			// handle special case
			//applet:IF_EGREP(APPLET_ODDNAME(egrep, grep, BB_DIR_BIN, BB_SUID_DROP, egrep))
			//applet:IF_FGREP(APPLET_ODDNAME(fgrep, grep, BB_DIR_BIN, BB_SUID_DROP, fgrep))
			if original == "egrep" || original == "fgrep" {
				new, ok := mapping["grep"]
				if ok {
					suffix = strings.ReplaceAll(suffix, ", grep,", fmt.Sprintf(", %s,", new))
				}
			}

			if newName, exists := mapping[original]; exists {
				// also replace parameter names in suffix
				var r2 = regexp.MustCompile(fmt.Sprintf(`,\s*%s`, original))
				suffix = r2.ReplaceAllString(suffix, fmt.Sprintf(", %s", newName))

				modifiedLine := prefix + newName + suffix
				fmt.Printf("%s: %s\n", path, modifiedLine)
				return modifiedLine
			}
			fmt.Printf("%s: %s\n", path, prefix+original+suffix)
			return prefix + original + suffix // return maybe modified suffix
		}
		return match
	})
}
