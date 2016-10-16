package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	// Default length of pomodoro
	defaultPomodoro time.Duration = time.Minute * 25
	// Default length of rest period
	defaultRest time.Duration = time.Minute * 5

	lenProgressBar int = 78
)

func clearTerminal() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[H")
}

func updateDisplay(remainingTime time.Duration, totalTime time.Duration, iterationNum int,
	name string) {
	/* Updates the terminal with time remaining and a progress bar */
	var minutes int = int(remainingTime.Minutes())
	var seconds int = int(math.Floor(remainingTime.Seconds()+0.5)) % 60
	var percentComplete float64 = float64(totalTime-remainingTime) / float64(totalTime)
	var progress int = int(percentComplete * float64(lenProgressBar))
	clearTerminal()
	fmt.Printf("%s %d\n", name, iterationNum)
	fmt.Printf("%v:%v (%v%%)\n", minutes, seconds, math.Floor(percentComplete*100))
	fmt.Printf("|%s%s|\n", strings.Repeat("-", progress),
		strings.Repeat(" ", lenProgressBar-progress))
}

func doIteration(duration time.Duration, iterationNum int, name string) {
	/* Runs a single iteration of a pomodoro or rest period. */
	var startTime time.Time = time.Now()
	ticker := time.NewTicker(time.Second)
	updateDisplay(duration, duration, iterationNum, name)
	for {
		<-ticker.C
		updateDisplay(duration-time.Since(startTime), duration, iterationNum, name)
		if time.Since(startTime) > duration {
			ticker.Stop()
			updateDisplay(duration-time.Since(startTime), duration, iterationNum, name)
			return
		}
	}
}

func printTransition(message string) {
	/* Prints a transition message */
	fmt.Printf("%s", message)
	for i := 0; i <= 5; i++ {
		time.Sleep(time.Second / 4)
		fmt.Printf(".")
	}
}

func main() {
	var iterationNum int = 1
	// How many pomodoros to run. If 0, run indefinitely.
	var iterations int = 0
	for {
		// Pomodoro
		doIteration(defaultPomodoro, iterationNum, "Pomodoro")
		printTransition(fmt.Sprintf("Pomodoro %d finished! Starting rest period",
			iterationNum))

		// Rest Period
		doIteration(defaultRest, iterationNum, "Rest Period")

		if iterationNum != 0 && iterationNum == iterations {
			if iterationNum == 1 {
				fmt.Println("Finished 1 Pomodoro!")
			} else {
				fmt.Printf("Finished %d Pomodoros!\n", iterationNum)
			}
			break
		} else {
			iterationNum += 1
			printTransition(fmt.Sprintf("Starting Pomodoro %d",
				iterationNum))
		}
	}
}
