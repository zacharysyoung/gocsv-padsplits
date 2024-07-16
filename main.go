// padsplits finds files created by gocsv's split subcommand and
// renames them by padding the numbers so that they can be sorted
// numerically.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: padsplits PREFIX

Find CSV files created by gocsv split, prefixed with PREFIX, like
out-1.csv...out-11.csv, and pad the numbers so the sort numerically,
like out-01.csv...out-11.csv.`)
		os.Exit(2)
	}
	flag.Parse()

	if n := len(flag.Args()); n != 1 {
		fatalf("got %d args; need one prefix", n)
	}

	prefix := strings.TrimSpace(flag.Arg(0))
	if prefix == "" {
		fatalf("got empty prefix; need one prefix", nil)
	}

	padSplits(prefix)

	// movedFiles := padSplits(prefix)
	// for _, file := range movedFiles {
	// 	fmt.Println(file.oldName, "→", file.newName)
	// }
}

func padSplits(prefix string) []file {
	if !strings.HasSuffix(prefix, "-") {
		prefix += "-"
	}

	globPattern := prefix + "*"
	paths, err := filepath.Glob(globPattern)
	if err != nil {
		fatalf("could not glob for %s: %v", globPattern, err)
	}

	files, maxLen, err := getFiles(paths, prefix)
	if err != nil {
		fatalf("could not get files: %v", err)
	}

	for i, file := range files {
		newName := fmt.Sprintf("%s%0*d.csv", prefix, maxLen, file.ord)
		err := move(file.oldName, newName)
		if err != nil {
			fatalf("could not move %s to %s: %v", file.oldName, newName, err)
		}
		files[i].newName = newName
	}

	return files
}

type file struct {
	ord     int
	oldName string
	newName string
}

func getFiles(paths []string, prefix string) (files []file, maxLen int, err error) {
	files = make([]file, 0)

	for _, name := range paths {
		x := strings.TrimPrefix(name, prefix)
		x = strings.TrimSuffix(x, ".csv")
		i, err := strconv.Atoi(x)
		if err != nil {
			return files, maxLen, fmt.Errorf("could not parse number in filename %s: %v", x, err)
		}

		files = append(files, file{ord: i, oldName: name})

		if n := len(x); n > maxLen {
			maxLen = n
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ord < files[j].ord
	})

	return files, maxLen, nil
}

func move(oldName, newName string) error {
	if oldName == newName {
		return nil
	}

	fIn, err := os.Open(oldName)
	if err != nil {
		return err
	}
	defer fIn.Close()

	fOut, err := os.Create(newName)
	if err != nil {
		return err
	}
	defer fOut.Close()

	if _, err := io.Copy(fOut, fIn); err != nil {
		return err
	}

	return os.Remove(oldName)
}

func fatalf(format string, args ...any) {
	if !strings.HasPrefix(format, "error: ") {
		format = "error: " + format
	}
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	if len(args) == 1 && args[0] == nil {
		args = nil
	}
	switch args {
	default:
		fmt.Fprintf(os.Stderr, format, args...)
	case nil:
		fmt.Fprintf(os.Stderr, format)
	}
	os.Exit(2)
}
