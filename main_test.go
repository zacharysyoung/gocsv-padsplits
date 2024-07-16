package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestMove(t *testing.T) {
	const data = "foo bar baz"

	var (
		tmpdir = t.TempDir()
		filea  = filepath.Join(tmpdir, "a")
		fileb  = filepath.Join(tmpdir, "b")
	)

	if err := write(filea, data); err != nil {
		t.Fatalf("could not write filea: %v", err)
	}
	if err := move(filea, fileb); err != nil {
		t.Fatalf("could not move filea to fileb: %v", err)
	}
	bdata, err := read(fileb)
	if err != nil {
		t.Fatalf("could not read fileb: %v", err)
	}

	if bdata != data {
		t.Errorf("bdata = %q; want %q", bdata, data)
	}

	names, err := getEntryNames(tmpdir)
	if err != nil {
		t.Fatalf("could not get entry names: %v", err)
	}

	if len(names) != 1 || names[0] != filepath.Base(fileb) {
		t.Errorf("got the following names for tmpdir: %v; want %s", names, fileb)
	}
}

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

// fileData concatenates the header line, "I", and the string-data
// for the file at fileOrd.
func fileData(fileOrd int) string {
	return "I\n" + files[fileOrd]
}

// TestPadSplits tests the end-to-end process: mimicking the
// creation of 11 CSV files (as if from gocsv split), calling
// padSplits, and verifing that only the correct padded names
// with their correct data exist.
func TestPadSplits(t *testing.T) {
	tmpdir := t.TempDir()

	for _, num := range ords {
		name := fmt.Sprintf("input-%d.csv", num)
		write(filepath.Join(tmpdir, name), fileData(num))
	}

	padSplits(filepath.Join(tmpdir, "input-"))

	refNames, err := getEntryNames(tmpdir)
	if err != nil {
		t.Fatalf("could not get list of names in tmpdir: %v", err)
	}

	if len(refNames) != 11 {
		t.Fatalf("after calling padSplits got %d entries in tmpdir, %v; want 11", len(refNames), refNames)
	}

	for _, num := range ords {
		name := fmt.Sprintf("input-%02d.csv", num) // hard-code %02d because only 11 files
		path := filepath.Join(tmpdir, name)

		got, err := read(path)
		if err != nil {
			t.Errorf("could not read file %s amongst %v: %v", name, refNames, err)
			continue
		}

		if want := fileData(num); got != want {
			t.Errorf("for file %s, got %s; want %s", name, got, want)
		}
	}
}

func write(path, data string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	f.WriteString(data)
	return f.Close()
}

func read(path string) (data string, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", errors.New("no such file or directory exists")
		} else {
			return "", err
		}
	}

	return string(b), nil
}

// getEntryNames returns the base names rooted at tmpdir.
func getEntryNames(tmpdir string) ([]string, error) {
	names := make([]string, 0)

	entries, err := os.ReadDir(tmpdir)
	if err != nil {
		return names, err
	}
	for _, x := range entries {
		names = append(names, filepath.Base(x.Name()))
	}

	return names, nil
}
