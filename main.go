package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

var _ = exec.Command

func isLong(d time.Duration) bool {
	return (d > 600*time.Millisecond && d < 1000*time.Millisecond )
}

func isShort(d time.Duration) bool {
	return (d < 600*time.Millisecond && d > 100*time.Millisecond)
}

func check(stack []time.Time) bool {
	a, b, c := stack[0], stack[1], stack[2]
	ps, os := a.Sub(b), b.Sub(c)
	fmt.Println(ps)
	return isShort(ps) && isShort(os)
}

func push(stack []time.Time, t time.Time) {
	for i := len(stack)-2; i>=0; i-- {
		stack[i+1] = stack[i]
	}
	stack[0] = t
}

func main() {
	var state bool
	stack := make([]time.Time, 3) // actually, a queue
	bf := bufio.NewReaderSize(os.Stdin, 32000)
	b := make([]byte, 16000) // one or two seconds

	for now := range time.Tick(20 * time.Millisecond) { // 50Hz...
		max, min := 0, 255

		n, err := bf.Read(b)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := 0; i < n; i++ {
			if int(b[i]) > max { max = int(b[i]) }
			if int(b[i]) < min { min = int(b[i]) }
		}

		amp := max-min

		if amp > 150 && !state{ //enforce pause
			state = true
			push(stack, now)
			if check(stack) {
				fmt.Println("\n**BEAT**\n")
				if len(os.Args) == 2 {
					go func() {
						cmd := exec.Command(os.Args[1])
						out, err := cmd.Output()
						if err != nil { fmt.Println(err) }
						fmt.Println(string(out))
					}()
				}
			}
		} else {
			state = false
		}
	}

}
