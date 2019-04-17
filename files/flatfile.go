package files

import (
	"bufio"
	"fl-library/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
TODO: This file needs cleaned up a bit
*/

type FlatFile struct {
	Name   string
	Labels []string
	Values [][]string
}

func unpackLine(s string) (output []string) {
	output = strings.Split(strings.Trim(s, "\n"), "^")
	return output[:len(output)-1]
}

func packLine(l []string) (output string) {
	l = append(l, "\n")
	return strings.Join(l, "^")
}

func (t *FlatFile) LoadFromFile(path string) {

	// Grab absolute path
	path, err := filepath.Abs(path)
	utils.Check(err)

	inFile, err := os.Open(path)
	utils.Check(err)
	defer inFile.Close()

	t.Name = strings.TrimSuffix(filepath.Base(inFile.Name()), t.getExt())

	r := bufio.NewReader(inFile)
	x, err := r.ReadByte()
	utils.Check(err)
	if x == '#' {
		l, _, err := r.ReadLine()
		utils.Check(err)
		t.Labels = unpackLine(string(l))

		for {
			line, _, err := r.ReadLine()
			if err != nil {
				break
			}
			t.Values = append(t.Values, unpackLine(string(line)))
		}
	} else {
		log.Fatal("Not a valid text map")
	}
}

func (t *FlatFile) DumpToFile(outdir string) {

	// Grab absolute path
	outdir, err := filepath.Abs(outdir)
	utils.Check(err)

	// Create directory
	err = os.MkdirAll(outdir, 0755)
	utils.Check(err)
	outdir += string(filepath.Separator) + t.Name + t.getExt()

	// Open out file
	outFile, err := os.OpenFile(outdir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	utils.Check(err)
	defer outFile.Close()

	utils.FileWrite(outFile, []byte(fmt.Sprintf("#%s", strings.ToUpper(packLine(t.Labels)))))
	for _, line := range t.Values {
		utils.FileWrite(outFile, []byte(packLine(line)))
	}
}

func (t *FlatFile) WriteToFile(outDir string) {

	// Grab absolute path
	outDir, err := filepath.Abs(outDir)
	utils.Check(err)

	// Create directory
	err = os.MkdirAll(outDir, 0755)
	utils.Check(err)
	outDir += string(filepath.Separator) + t.Name + "_CLEAN" + t.getExt()

	// Open outfile
	outFile, err := os.OpenFile(outDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	utils.Check(err)
	defer outFile.Close()

	var objs []string
	for _, line := range t.Values {
		var obj []string
		for i, l := range t.Labels {
			obj = append(obj, fmt.Sprintf("%s=%s", l, line[i]))
		}
		objs = append(objs, strings.Join(obj, ",\n"))
	}
	utils.FileWriteString(outFile, strings.Join(objs, "\n::\n"))
}

func (t *FlatFile) LoadFromCleanFile(path string) {

	// Grab absolute path
	path, err := filepath.Abs(path)
	utils.Check(err)

	inFile, err := os.Open(path)
	utils.Check(err)
	defer inFile.Close()

	t.Name = strings.TrimSuffix(filepath.Base(inFile.Name()), "_CLEAN"+t.getExt())
	t.Name += "_MOD"

	fs, err := inFile.Stat()
	utils.Check(err)
	data := make([]byte, int64(fs.Size()))
	utils.FileRead(inFile, data)

	lines := strings.Split(string(data), "\n::\n")

	for _, o := range lines {
		var obj []string
		fields := strings.Split(o, ",\n")
		for _, f := range fields {
			pair := strings.Split(f, "=")
			if !utils.StringInSlice(pair[0], t.Labels) {
				t.Labels = append(t.Labels, pair[0])
			}
			obj = append(obj, pair[1])
		}
		t.Values = append(t.Values, obj)
	}
}

func (f *FlatFile) getExt() string {
	return ".txt"
}
