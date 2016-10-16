package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
		time.Sleep(time.Second / 5)
		fmt.Printf(".")
	}
}

func printSummary(finishedPomodoros int, startTime time.Time) {
	/* Prints a summary of the work completed */
	if finishedPomodoros == 1 {
		fmt.Println("\nFinished 1 Pomodoro!")
	} else {
		fmt.Printf("\nFinished %d Pomodoros!\n", finishedPomodoros)
	}
	fmt.Printf("Elapsed Time: %v\n", time.Since(startTime))
}

func main() {
	var iterationNum int = 1
	// How many pomodoros to run. If 0, run indefinitely.
	var iterations int = 0
	var resting bool = false
	var startTime time.Time = time.Now()

	// Cleanup on CTRL-C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig
		var finishedPomodoros int = iterationNum
		if !resting {
			finishedPomodoros = iterationNum - 1
		}
		printSummary(finishedPomodoros, startTime)
		os.Exit(0)
	}()

	for {
		// Pomodoro
		doIteration(defaultPomodoro, iterationNum, "Pomodoro")
		resting = true
		printTransition(fmt.Sprintf("Pomodoro %d finished! Starting rest period",
			iterationNum))

		// Rest Period
		doIteration(defaultRest, iterationNum, "Rest Period")

		if iterationNum != 0 && iterationNum == iterations {
			printSummary(iterationNum, startTime)
			break
		} else {
			iterationNum += 1
			resting = false
			printTransition(fmt.Sprintf("Starting Pomodoro %d",
				iterationNum))
		}
	}
}
