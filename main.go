package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	out, err := exec.Command("git", "log", "--format='%at'").Output()
	if err != nil {
		log.Fatalln(string(out), err)
	}

	split := strings.Fields(string(out))
	for i := range split {
		split[i] = strings.TrimSpace(split[i])
		split[i] = strings.Trim(split[i], "'")
		split[i] = strings.TrimSpace(split[i])
	}

	times := make([]time.Time, len(split))
	for i := range split {
		sec, err := strconv.ParseInt(split[i], 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		times[i] = time.Unix(sec, 0)
		if err != nil {
			log.Fatalln(err)
		}
	}

	perHour := make(map[int]int)

	for _, t := range times {
		perHour[t.Local().Hour()]++
	}

	min := math.MaxInt
	max := -1
	for _, cnt := range perHour {
		if cnt < min {
			min = cnt
		}
		if cnt > max {
			max = cnt
		}
	}

	log10 := int(math.Ceil(math.Log10(float64(max))))

	termWidth, err := getWidth()
	if err != nil {
		log.Fatalln(err)
	}

	scaleFactor := math.Min(80, float64(termWidth)-float64(8+log10)) / float64(max)

	for i := 0; i < 24; i++ {
		fmt.Printf(fmt.Sprintf("%%02d: %%%dd  %%s\n", log10), i, perHour[i], strings.Repeat("#", int(float64(perHour[i])*scaleFactor)))
	}
}

func getWidth() (int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	splitted := strings.Split(strings.TrimSpace(string(out)), " ")
	width := splitted[1]

	return strconv.Atoi(string(width))
}
