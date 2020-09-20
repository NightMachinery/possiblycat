/// 2>/dev/null ; exec gorun "$0" "$@"

package main

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {
	wait := 10
	if len(os.Args) >= 2 {
		waitDummy, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		wait = waitDummy
	}

	// disable input buffering
	// exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	// exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	done := make(chan bool, 1)
	b := make(chan byte, 1)
	go scan(b, done)

	select {
	case res := <-b:
		inBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		stdin := append([]byte{res}, inBytes...)
		_, err2 := os.Stdout.Write(stdin)
		if err2 != nil {
			panic(err2)
		}
	case <-done:
		os.Exit(0)
	case <-time.After(time.Duration(wait) * time.Millisecond):
		os.Exit(1)
	}
}

func scan(out chan byte, done chan bool) {
	var b []byte = make([]byte, 1)
	_, err := os.Stdin.Read(b)
	if err == io.EOF {
		done <- true
		return
	} else if err != nil {
		panic(err)
	}
	out <- b[0]
}
