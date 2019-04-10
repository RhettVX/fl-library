package files

import (
	"fl-library/utils"
	"os"
	"path/filepath"
)

type LocaleEntry struct {
	Id   int64
	Flag string
	Text string
}

type Locale struct {
	Name    string
	Entires []LocaleEntry
}

func (l *Locale) LoadFromDir(path string) {

	// Grab absolute path
	path, err := filepath.Abs(path)
	utils.Check(err)

	var datPath, dirPath string
	// Walk through dir to find english files
	err = filepath.Walk(path, func(file string, _ os.FileInfo, e error) error {
		basename := filepath.Base(file)
		if basename == "en_us_data.dat" {
			datPath = file
		} else if basename == "en_us_data.dir" {
			dirPath = file
		}
		return e
	})
	utils.Check(err)
	// TODO: LEFT OFF HERE
}

func (l *Locale) getId() []byte {
	return []byte{0x3f, 0xbb, 0xbf}
}
