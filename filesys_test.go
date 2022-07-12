package main

import (
	"fmt"
	"os"
	"testing"
)

func TestSkim(t *testing.T) {

	RecursiveTouch("../sonr/internal/motor", func(file *os.File) {
		filestat, _ := file.Stat()
		fmt.Println("==" + filestat.Name())
	})
}
