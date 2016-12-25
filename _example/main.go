package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sheenobu/sjson/pkg/sjson"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: sjsontest <filename>\n")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	ch := make(chan sjson.Token)
	go func() {
		defer close(ch)
		err = sjson.ReadAll(f, ch)
		if err != nil && err != io.EOF {
			panic(err)
		}
	}()

	var offset int
	for {
		select {
		case t, more := <-ch:
			if !more {
				return
			}

			if t.Type() == sjson.EndType {
				offset--
			}
			for i := 0; i <= offset; i++ {
				fmt.Printf(" ")
			}
			fmt.Printf("%s\n", t)

			if t.Type() == 5 || t.Type() == 6 {
				offset++
			}

		}
	}
}
