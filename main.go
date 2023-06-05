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

	perHour := make(map[int]int, 24)

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

	termWidth, err := getWidth()
	if err != nil {
		log.Fatalln(err)
	}

	log10 := int(math.Ceil(math.Log10(float64(max))))

	delimiter := "  "
	// every thing is double escaped because we need to determine the format string first.
	// afterwards the formatString looks something like this:
	// "%02d: %3d%s\n"
	//                                    vv------- this is set by log10
	formatString := fmt.Sprintf("%%02d: %%%dd%%s\n", log10)
	// minus 1 because of the newline
	prefixLen := len(fmt.Sprintf(formatString, 23, max, delimiter)) - 1

	// Max width should be 80
	width := math.Min(float64(termWidth-prefixLen), float64(80-prefixLen))

	// must have at least one hashtag
	scaleFactor := math.Max(width/float64(max), 0)

	for hour, amount := range perHour {
		numberOfHashtags := int(math.Round(float64(amount) * scaleFactor))
		hashtags := delimiter + strings.Repeat("#", numberOfHashtags)
		if scaleFactor == 0 {
			hashtags = ""
		}

		fmt.Printf(formatString, hour, amount, hashtags)
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
