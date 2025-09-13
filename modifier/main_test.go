package main

import (
	"regexp"
	"testing"
)

func TestRegexes(t *testing.T) {
	tests := []struct {
		input string
		regex *regexp.Regexp
		name  string
	}{
		{input: `//usage:#define tar_trivial_usage`, regex: regexUsage, name: "tar"},
		{input: `//usage:#define tar_full_usage "\n\n"`, regex: regexUsage, name: "tar"},
		{input: `//usage:#define tar_example_usage`, regex: regexUsage, name: "tar"},
		{input: `//usage:#define tree_trivial_usage NOUSAGE_STR`, regex: regexUsage, name: "tree"},
		{input: `//usage:#define tree_full_usage ""`, regex: regexUsage, name: "tree"},
		{input: `exit(ls_main(/*argc_unused*/ 0, (char**) argv));`, regex: regexMain, name: "ls"},
		{input: `int id_main(int argc, char **argv) MAIN_EXTERNALLY_VISIBLE;`, regex: regexMain, name: "id"},
		{input: `int id_main(int argc UNUSED_PARAM, char **argv)`, regex: regexMain, name: "id"},
		{input: `int blkid_main(int argc, char **argv) MAIN_EXTERNALLY_VISIBLE;`, regex: regexMain, name: "blkid"},
		{input: `int blkid_main(int argc UNUSED_PARAM, char **argv)`, regex: regexMain, name: "blkid"},
		{input: `ret = run_applet_main(argv, kill_main);`, regex: regexMain2, name: "kill"},
		{input: `	BLTIN("eval"     , builtin_eval    , "Construct and run shell command"),`, regex: regexBuiltIn, name: "eval"},
		{input: `	BLTIN("export"   , builtin_export  , "Set environment variables"),`, regex: regexBuiltIn, name: "export"},
		{input: `//applet:IF_VLOCK(APPLET(vlock, BB_DIR_USR_BIN, BB_SUID_REQUIRE))`, regex: regexApplet, name: "vlock"},
		{input: `//applet:IF_PWDX(APPLET_NOFORK(pwdx, pwdx, BB_DIR_USR_BIN, BB_SUID_DROP, pwdx))`, regex: regexApplet, name: "pwdx"},
		{input: `//applet:IF_RM(APPLET_NOEXEC(rm, rm, BB_DIR_BIN, BB_SUID_DROP, rm))`, regex: regexApplet, name: "rm"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			matches := test.regex.FindStringSubmatch(test.input)
			if len(matches) != 4 {
				t.Errorf("Expected 4 matches, got %d", len(matches))
			} else if matches[2] != test.name {
				t.Errorf("Expected name %s, got %s", test.name, matches[2])
			}
		})
	}
}
