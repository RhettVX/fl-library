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

type TextMap struct {
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

func (t *TextMap) LoadFromFile(path string) {

	// Grab absolute path
	path, err := filepath.Abs(path)
	utils.Check(err)

	inFile, err := os.Open(path)
	utils.Check(err)
	defer inFile.Close()

	t.Name = strings.TrimSuffix(filepath.Base(inFile.Name()), ".txt")

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

func (t *TextMap) DumpToFile(outdir string) {

	// Grab absolute path
	outdir, err := filepath.Abs(outdir)
	utils.Check(err)

	// Create directory
	err = os.MkdirAll(outdir, 0666)
	utils.Check(err)
	outdir += `\` + t.Name + ".txt"

	// Open out file
	outFile, err := os.OpenFile(outdir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	utils.Check(err)
	defer outFile.Close()

	outFile.Write([]byte(fmt.Sprintf("#%s", strings.ToUpper(packLine(t.Labels)))))
	for _, line := range t.Values {
		outFile.Write([]byte(packLine(line)))
	}
}

func (t *TextMap) WriteToFile(outDir string) {

	// Grab absolute path
	outDir, err := filepath.Abs(outDir)
	utils.Check(err)

	// Create directory
	err = os.MkdirAll(outDir, 0666)
	utils.Check(err)
	outDir += `\` + t.Name + "_clean.txt"

	// Open outfile
	outFile, err := os.OpenFile(outDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
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
	outFile.WriteString(strings.Join(objs, "\n::\n"))
}

func (t *TextMap) LoadFromCleanFile(path string) {

	// Grab absolute path
	path, err := filepath.Abs(path)
	utils.Check(err)

	inFile, err := os.Open(path)
	utils.Check(err)
	defer inFile.Close()

	t.Name = strings.TrimSuffix(filepath.Base(inFile.Name()), "_clean.txt")
	t.Name += "_MOD"

	fs, err := inFile.Stat()
	utils.Check(err)
	data := make([]byte, int64(fs.Size()))
	inFile.Read(data)

	utils.Check(err)
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
