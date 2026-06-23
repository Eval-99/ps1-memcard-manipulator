package main

import (
	"errors"
	"fmt"
	"strings"
)

var SUPPORTED_FORMATS = [5]string{".mcr", ".bin", ".gme", ".mcd", ".srm"}

const (
	BLOCK_SIZE       = 8192
	CARD_SIZE        = 131072
	DIR_FRAME_OFFSET = 128
	MCS_FRAME_SIZE   = 128 + BLOCK_SIZE
)

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

func main() {
	filename := "somefile.exe"

	result, err := getFileExtension(filename)
	if err != nil {
		fmt.Println(err)
	}

	thing, err := changeFileExtension(filename, "zip")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(filename, result, "=======", thing)
}
