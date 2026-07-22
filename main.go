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

func updateChecksum(data [CARD_SIZE]byte, offset int) {
	var x byte = 0
	for k := range 127 {
		x ^= data[offset+k]
	}
	data[offset+127] = x
}

func verifyChecksum(data [CARD_SIZE]byte, offset int) bool {
	var x byte = 0
	for k := range 127 {
		x ^= data[offset+k]
	}
	return data[offset+127] == x
}

func parseString(data [CARD_SIZE]byte, offset int, length int) string {
	var s strings.Builder
	for k := range length {
		s.WriteString(string(data[offset+k]))
	}
	return s.String()
}

func parseShiftJIS(data [CARD_SIZE]byte, offset int, length int) (string, error) {
	sub := data[offset : offset+length]
	end := slices.Index(sub, 0)
	if end >= 0 {
		sub = sub[0:end]
	}
	decoder := japanese.ShiftJIS.NewDecoder()
	return decoder.String(string(sub))
}

// --- Card Structure ---

func getLinkedBlocks(data [CARD_SIZE]byte, startSlot int) []int {
	slots := []int{startSlot}
	currentSlot := startSlot

	for range 15 {
		entryOffset := DIR_FRAME_OFFSET + (currentSlot * 128)
		linkVal := uint16(data[entryOffset+8]) | (uint16(data[entryOffset+9]) << 8)

		if linkVal == 0xFFFF {
			break
		}
		if linkVal >= 15 {
			break
		}
		if slices.Contains(slots, int(linkVal)) {
			break
		}

		slots = append(slots, int(linkVal))
		currentSlot = int(linkVal)
	}

	return slots
}

func findFreeSlots(data [CARD_SIZE]byte, count int) []int {
	freeSlots := []int{}

	for i := range 15 {
		status := data[DIR_FRAME_OFFSET+(i*128)]
		if status == 0xA0 {
			freeSlots = append(freeSlots, i)
		}
	}

	if len(freeSlots) < count {
		return nil
	}

	return freeSlots[0:count]
}

func slotHasData(data [CARD_SIZE]byte, slotIndex int) bool {
	blockStart := (slotIndex + 1) * BLOCK_SIZE

	for i := range 128 {
		if data[blockStart+i] != 0 {
			return true
		}
	}

	return false
}

func countUsedBlocks(data [CARD_SIZE]byte) int {
	used := 0

	for i := range 15 {
		status := data[DIR_FRAME_OFFSET+(i*128)]
		if status == 0x51 || status == 0x52 || status == 0x53 {
			used++
		}
	}

	return used
}

func getSlotStatus(data [CARD_SIZE]byte, slotIndex int) byte {
	return data[DIR_FRAME_OFFSET+(slotIndex*128)]
}

// --- Card Operations ---

func createBlankCard() [CARD_SIZE]byte {
	data := [CARD_SIZE]byte{}

	data[0] = 0x4D // 'M'
	data[1] = 0x43 // 'C'
	data[127] = 0x0E

	for i := range 15 {
		off := DIR_FRAME_OFFSET + (i * 128)
		data[off] = 0xA0
		data[off+4] = 0
		data[off+5] = 0
		updateChecksum(data, off)
	}

	return data
}

func main() {
	data := [CARD_SIZE]byte{
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
