package main

import (
	"fmt"
)

type screen struct {
	w, h  int
	chars []rune
}

func (s screen) set(y, x int, c rune) {
	s.chars[y*s.w+x] = c
}

func (s screen) clear() {
	for i := range s.chars {
		s.chars[i] = ' '
	}
}

func (s screen) String() string {
	st := ""
	for i := 0; i < s.h; i++ {
		st += fmt.Sprintln(string(s.chars[i*s.w : (i+1)*s.w]))
	}
	return st
}
