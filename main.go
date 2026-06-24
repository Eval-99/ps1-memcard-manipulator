package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

var SUPPORTED_FORMATS = [5]string{".mcr", ".bin", ".gme", ".mcd", ".srm"}

const (
	BLOCK_SIZE       = 8192
	CARD_SIZE        = 131072
	DIR_FRAME_OFFSET = 128
	MCS_FRAME_SIZE   = 128 + BLOCK_SIZE
)

// --- Utility ---

func getFileExtension(filename string) (string, error) {
	if len(filename) == 0 {
		return "", errors.New("Error, no file.")
	}

	strSlice := strings.Split(filename, ".")
	if len(strSlice) == 1 {
		return "", errors.New("Error, no current file extension.")
	}

	return strings.ToLower(strSlice[len(strSlice)-1]), nil
}

func changeFileExtension(filename, newExt string) (string, error) {
	if len(filename) == 0 || len(newExt) == 0 {
		return "", errors.New("Error, no file or no new file extension.")
	}

	strSlice := strings.Split(filename, ".")
	if len(strSlice) == 1 {
		return "", errors.New("Error, no current file extension.")
	}

	strSlice = append(strSlice[0:len(strSlice)-1], newExt)

	return strings.ToLower(strings.Join(strSlice, ".")), nil
}

// --- Binary Parsing ---

func updateChecksum(data []byte, offset int) {
	var x byte = 0
	for k := range 127 {
		x ^= data[offset+k]
	}
	data[offset+127] = x
}

func verifyChecksum(data []byte, offset int) bool {
	var x byte = 0
	for k := range 127 {
		x ^= data[offset+k]
	}
	return data[offset+127] == x
}

func parseString(data []byte, offset int, length int) string {
	var s strings.Builder
	for k := range length {
		s.WriteString(string(data[offset+k]))
	}
	return s.String()
}

func parseShiftJIS(data []byte, offset int, length int) (string, error) {
	sub := data[offset : offset+length]
	end := slices.Index(sub, 0)
	if end >= 0 {
		sub = sub[0:end]
	}
	decoder := japanese.ShiftJIS.NewDecoder()
	return decoder.String(string(sub))
}

func main() {
	data := []byte{
		130, 177,
		130, 241,
		130, 201,
		130, 191,
		130, 205,
		129, 65,
		130, 168,
		140, 179,
		139, 67,
		130, 197,
		130, 183,
		130, 169,
		129, 72,
	}

	japaneseString, _ := parseShiftJIS(data, 0, 26)
	fmt.Println(japaneseString)
}
