package main

import (
	"fmt"
	"github.com/gosuri/uilive"
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

func updateDisplay(writer *uilive.Writer, remainingTime time.Duration,
	totalTime time.Duration, iterationNum int, name string) {
	/* Updates the terminal with time remaining and a progress bar */
	var minutes int = int(remainingTime.Minutes())
	var seconds int = int(math.Floor(remainingTime.Seconds()+0.5)) % 60
	var percentComplete float64 = float64(totalTime-remainingTime) / float64(totalTime)
	var progress int = int(percentComplete * float64(lenProgressBar))
	fmt.Fprintf(writer, "%s %d\n", name, iterationNum)
	fmt.Fprintf(writer, "%v:%v (%v%%)\n", minutes, seconds, math.Floor(percentComplete*100))
	fmt.Fprintf(writer, "|%s%s|\n", strings.Repeat("-", progress),
		strings.Repeat(" ", lenProgressBar-progress))
	writer.Flush()
}

func doIteration(writer *uilive.Writer, duration time.Duration,
	iterationNum int, name string) {
	/* Runs a single iteration of a pomodoro or rest period. */
	var startTime time.Time = time.Now()
	ticker := time.NewTicker(time.Second)
	updateDisplay(writer, duration, duration, iterationNum, name)
	for {
		<-ticker.C
		updateDisplay(writer, duration-time.Since(startTime), duration, iterationNum, name)
		if time.Since(startTime) > duration {
			ticker.Stop()
			updateDisplay(writer, duration-time.Since(startTime), duration, iterationNum, name)
			return
		}
	}
}

func printTransition(writer *uilive.Writer, message string) {
	/* Prints a transition message */
	fmt.Fprintf(writer, "%s", message)
	for i := 0; i <= 5; i++ {
		time.Sleep(time.Second / 5)
		fmt.Fprintf(writer, ".")
	}
	fmt.Fprintf(writer, "\n")
	writer.Flush()
}

func printSummary(writer *uilive.Writer, finishedPomodoros int, startTime time.Time) {
	/* Prints a summary of the work completed */
	if finishedPomodoros == 1 {
		fmt.Fprintf(writer, "\nFinished 1 Pomodoro!\n")
	} else {
		fmt.Fprintf(writer, "\nFinished %d Pomodoros!\n", finishedPomodoros)
	}
	fmt.Fprintf(writer, "Elapsed Time: %v\n", time.Since(startTime))
}

func main() {
	var iterationNum int = 1
	// How many pomodoros to run. If 0, run indefinitely.
	var iterations int = 0
	var resting bool = false
	var startTime time.Time = time.Now()
	writer := uilive.New()
	writer.Start()

	// Cleanup on CTRL-C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig
		var finishedPomodoros int = iterationNum
		if !resting {
			finishedPomodoros = iterationNum - 1
		}
		printSummary(writer, finishedPomodoros, startTime)
		writer.Stop()
		os.Exit(0)
	}()

	for {
		// Pomodoro
		doIteration(writer, defaultPomodoro, iterationNum, "Pomodoro")
		resting = true
		printTransition(writer, fmt.Sprintf("Pomodoro %d finished! Starting rest period",
			iterationNum))

		// Rest Period
		doIteration(writer, defaultRest, iterationNum, "Rest Period")

		if iterationNum != 0 && iterationNum == iterations {
			printSummary(writer, iterationNum, startTime)
			break
		} else {
			iterationNum += 1
			resting = false
			printTransition(writer, fmt.Sprintf("Starting Pomodoro %d",
				iterationNum))
		}
	}
	writer.Stop()
}
