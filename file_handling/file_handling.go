package file_handling

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const bytesPerLine = 32

func GetLastLines(numLines int, fileName string) ([]string, error) {
	// Open file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the size of the file
	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}

	// To read the last lines of the file, we need to guess a place near the end, such that we can barely fill up a
	// buffer with file data before we reach EOF. If this procedure doesn't result in numLines lines, we try again with
	// a bigger buffer.
	for multiplier := 1; ; multiplier++ {
		offset := multiplier * numLines * bytesPerLine
		buf := make([]byte, offset)
		start := stat.Size() - int64(offset) // We'll start reading near the end of the file with enough offset to fill the buffer
		if start < 0 {
			start = 0                       // Can't start at positions before 0
			buf = make([]byte, stat.Size()) // Make buffer for entire file
		}

		numBytes, err := file.ReadAt(buf, start)
		if err != nil {
			return nil, err
		}

		text := string(buf[:numBytes])
		lines := strings.Split(text, "\n")

		if start != 0 {
			lines = lines[1:] // Discard the first string as it might be clipped
		}

		if lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1] // Remove trailing empty string
		}

		if len(lines) >= numLines {
			return lines[len(lines)-numLines:], nil
		}

		if start == 0 {
			return lines, nil // Stop trying to get more lines if we reach the beginning
		}
		multiplier++
	}
	return nil, errors.New("unknown error")
}

func WriteLine(line string, fileName string) error {
	// Open file
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", line))
	return err
}
