package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "quiz.csv", "plik csv z pytaniami ")
	timeLimit := flag.Int("limit", 30, "maksymalny czas na wykonanie quizu")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Nie można otworzyć pliku: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Zły format pliku.")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	start := time.Now()
problemloop:
	for i, p := range problems {
		
		fmt.Printf("Pytanie #%d: %s ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\n Niestety twój czas się skończył. Na wykonanie quizu miałeś 30s")
			break problemloop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	stop:= time.Since(start)
	fmt.Printf("Zdobyłeś %d z %d możliwych do zdobycia punktów.\n", correct, len(problems))
	fmt.Printf("Twój czas %s\n", stop)
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}