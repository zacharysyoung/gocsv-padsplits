package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

var (
	files = map[int]string{
		1:  "a\n",
		2:  "b\n",
		3:  "c\n",
		4:  "d\n",
		5:  "e\n",
		6:  "f\n",
		7:  "g\n",
		8:  "h\n",
		9:  "i\n",
		10: "j\n",
		11: "k\n",
	}
	ords = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
)

// data concatenates the header line, "I", and the string-data
// for the file at fileOrd.
func data(fileOrd int) string {
	return "I\n" + files[fileOrd]
}

func TestPadSplits(t *testing.T) {
	tmpdir := t.TempDir()

	for _, num := range ords {
		name := fmt.Sprintf("input-%d.csv", num)
		f, err := os.Create(filepath.Join(tmpdir, name))
		if err != nil {
			t.Fatal(err)
		}
		f.WriteString(data(num))
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
	}

	padSplits(filepath.Join(tmpdir, "input-"))

	// Get reference list of files after padSplits, for error reporting
	refEntries, err := os.ReadDir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	refNames := make([]string, 0)
	for _, x := range refEntries {
		refNames = append(refNames, x.Name())
	}

	for _, num := range ords {
		name := fmt.Sprintf("input-%02d.csv", num) // hard-code %02d because only 11 files
		path := filepath.Join(tmpdir, name)

		b, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				t.Errorf("tried to open %s amongst %v: no such file or directory", name, refNames)
			} else {
				t.Errorf("tried to open %s amongst %v: %v", name, refNames, err)
			}
			continue
		}

		got := string(b)
		want := data(num)

		if got != want {
			t.Errorf("for %s, got %s; want %s", name, got, want)
		}
	}
}
