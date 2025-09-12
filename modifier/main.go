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

var regexApplet = regexp.MustCompile(`^(//applet:.+\(\s*APPLET[_A-Z]*\()([^,]+)(, .*)$`)

/*
//usage:#define tar_trivial_usage
//usage:#define tar_full_usage "\n\n"
//usage:#define tar_example_usage
//usage:#define tree_trivial_usage NOUSAGE_STR
//usage:#define tree_full_usage ""
*/
var regexUsage = regexp.MustCompile(`^(\s*//usage:#define )([a-z]+)(_[a-z]+_usage\s*.*)$`)

// exit(ls_main(/*argc_unused*/ 0, (char**) argv));
// int id_main(int argc, char **argv) MAIN_EXTERNALLY_VISIBLE;
// int id_main(int argc UNUSED_PARAM, char **argv)
// int blkid_main(int argc, char **argv) MAIN_EXTERNALLY_VISIBLE;
// int blkid_main(int argc UNUSED_PARAM, char **argv)
var regexMain = regexp.MustCompile(`^(.+?)([a-z]+)(_main\s*\(.+)`)

var mapping = map[string]string{
	"ls":     "zoag",
	"tar":    "asphalt",
	"man":    "aunleitung",
	"unzip":  "auspockn",
	"sed":    "austausch",
	"id":     "ausweis",
	"groups": "gruppn",
	"tree":   "bam",
	"vi":     "bittelassmichgehenichkommnichtraus",
	"mv":     "beweg",
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

	// skip files that don't contain the target string
	if !bytes.Contains(contents, []byte("//applet:")) {
		return nil
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
		case regexMain.MatchString(t):
			modified := process(path, t, regexMain)
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

			if newName, exists := mapping[original]; exists {
				// also replace parameter names in suffix
				var r2 = regexp.MustCompile(fmt.Sprintf(`,\s*%s`, original))
				suffix = r2.ReplaceAllString(suffix, fmt.Sprintf(", %s", newName))

				// handle special case of
				//applet:IF_GROUPS(APPLET_NOEXEC(groups, id, BB_DIR_USR_BIN, BB_SUID_DROP, groups))
				//applet:IF_ID(    APPLET_NOEXEC(id,     id, BB_DIR_USR_BIN, BB_SUID_DROP, id    ))
				// where groups reference id
				if original == "groups" {
					newId, ok := mapping["id"]
					if ok {
						suffix = strings.ReplaceAll(suffix, ", id,", fmt.Sprintf(", %s,", newId))
					}
				}

				modifiedLine := prefix + newName + suffix
				fmt.Printf("%s: %s\n", path, modifiedLine)
				return modifiedLine
			}
		}
		return match
	})
}
