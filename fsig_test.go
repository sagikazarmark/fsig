package main

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func TestFsig(t *testing.T) {
	os.Args = []string{"fsig", "-w", "test", "HUP", "--", "sh", "test.sh"}

	go func() {
		time.Sleep(2 * time.Second)

		os.Create("test/"+strconv.FormatInt(time.Now().Unix(), 10))
	}()

	main()
}
