package file_handling

import (
	"fmt"
	"os"
	"testing"
)

const fileName = "testfile"

func check(err error, t *testing.T) {
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

// This test assumes there exists a file with at least numExpectedLines lines of text in it
func TestGetLastLines(t *testing.T) {
	/* Setup */
	file, err := os.Create(fileName)
	check(err, t)
	file.Close()
	defer os.Remove(fileName)

	err = WriteLine("Potentially clipped string", fileName)
	check(err, t)
	err = WriteLine("Hello world", fileName)
	check(err, t)
	err = WriteLine("Test test", fileName)
	check(err, t)
	err = WriteLine("Last line", fileName)
	check(err, t)

	/* Test */
	lines, err := GetLastLines(3, fileName)
	check(err, t)
	if len(lines) != 3 {
		fmt.Println("Expected to read 3 lines, got ", len(lines))
		t.Fail()
	}

	lines, err = GetLastLines(5, fileName)
	check(err, t)
	if len(lines) != 4 {
		fmt.Println("Expected to read 4 lines, got ", len(lines)) // Outputs the max amount of lines
		t.Fail()
	}

	lines, err = GetLastLines(1, fileName)
	check(err, t)
	if lines[0] != "Last line" {
		fmt.Println("Last line not as expected. Expected Last line, got", lines[0])
		t.Fail()
	}
}
