package main

import (
	"log"
	"os"
)

func init() {
	for i := range notes {
		if i == 0 {
			notes[0] = BOTTOM
			continue
		}
		notes[i] = RT12 * notes[i-1]

	}

	midnotes[0], midnotes[len(midnotes)-1] = 0.0, SAMP/2
	for i := 1; i < len(notes); i++ {
		midnotes[i] = (notes[i] + notes[i-1]) / 2
	}

	logFile, _ := os.Create("log")
	log.SetOutput(logFile)

	for i, v := range notes {
		log.Println(noteName(i), v, midnotes[i], midnotes[i+1])
	}
}

var notes = make([]float64, SEMI_N)      // 112 semitones between bottom and ~20kHz
var midnotes = make([]float64, SEMI_N+1) //boundary values; taking the halfway is somewhat arbitrary

var noteStrings = [...]string{
	"A",
	"A#",
	"B",
	"C",
	"C#",
	"D",
	"D#",
	"E",
	"F",
	"F#",
	"G",
	"G#",
}

func noteName(i int) string {
	return noteStrings[i%12]
}
