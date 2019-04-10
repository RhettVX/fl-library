package utils

import (
	"bufio"
	"hash/crc32"
	"log"
	"os"
	"strings"
)

func Check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func CalcCrc32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}

func TakeInput() (output string) {
	r := bufio.NewReader(os.Stdin)
	output, err := r.ReadString('\n')
	Check(err)

	output = strings.Trim(output, "\n\r\"")
	return output
}
