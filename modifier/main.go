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
)

var mapping = map[string]string{
	"ls": "zoag",
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
		if !info.IsDir() && filepath.Ext(info.Name()) == ".c" {
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
		if bytes.Contains([]byte(t), []byte("//applet:")) {
			modified := processApplet(t)
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

func processApplet(line string) string {
	fmt.Println("Processing applet line:", line)
	r := regexp.MustCompile(`^(//applet:.+\(APPLET[_A-Z]*\()([^,]+)(, .*)`)

	return r.ReplaceAllStringFunc(line, func(match string) string {
		matches := r.FindStringSubmatch(match)
		if len(matches) == 4 {
			prefix := matches[1]   // Everything before the name
			original := matches[2] // The original name (first capture group)
			suffix := matches[3]   // Everything after the name

			if newName, exists := mapping[original]; exists {
				modifiedLine := prefix + newName + suffix
				return modifiedLine
			}
		}
		return match
	})
}
