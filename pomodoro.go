package main

import (
	"flag"
	"fmt"
	"github.com/gosuri/uilive"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Command line options
type options struct {
	length     time.Duration // Length of pomodoro
	restLength time.Duration // Length of rest period
	iterations int           // Number of iterations before exiting
}

func updateDisplay(writer *uilive.Writer, remainingTime time.Duration,
	totalTime time.Duration, iterationNum int, name string) {
	const lenProgressBar int = 78
	/* Updates the terminal with time remaining and a progress bar */
	minutes := int(remainingTime.Minutes())
	seconds := int(math.Floor(remainingTime.Seconds()+0.5)) % 60
	percentComplete := float64(totalTime-remainingTime) / float64(totalTime)
	progress := int(percentComplete * float64(lenProgressBar))
	fmt.Fprintln(writer, name, iterationNum)
	fmt.Fprintf(writer, "%v:%v (%v%%)\n", minutes, seconds, math.Floor(percentComplete*100))
	fmt.Fprintf(writer, "|%s%s|\n", strings.Repeat("-", progress),
		strings.Repeat(" ", lenProgressBar-progress))
	writer.Flush()
}

func doIteration(writer *uilive.Writer, duration time.Duration,
	iterationNum int, name string) {
	/* Runs a single iteration of a pomodoro or rest period. */
	startTime := time.Now()
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
		fmt.Fprintln(writer, "\nFinished 1 Pomodoro!")
	} else {
		fmt.Fprintf(writer, "\nFinished %d Pomodoros!\n", finishedPomodoros)
	}
	fmt.Fprintln(writer, "Elapsed Time: ", time.Since(startTime))
}

func setupOptions() options {
	// Retrieves options from command line flags
	pomodoro := flag.Float64("length", 25, "Length of the pomodoro, in minutes")
	rest := flag.Float64("rest", 5, "Length of the rest period, in minutes")
	iterations := flag.Int("iterations", 0, "Number of iterations to run before exiting. If zero, run indefinetely")
	flag.Parse()
	return options{time.Second * time.Duration(*pomodoro*60), time.Second * time.Duration(*rest*60), *iterations}
}

func main() {
	opts := setupOptions()
	iterationNum := 1
	resting := false
	startTime := time.Now()
	writer := uilive.New()
	writer.Start()

	// Cleanup on CTRL-C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig
		finishedPomodoros := iterationNum
		if !resting {
			finishedPomodoros = iterationNum - 1
		}
		printSummary(writer, finishedPomodoros, startTime)
		writer.Stop()
		os.Exit(0)
	}()

	for {
		// Pomodoro
		doIteration(writer, opts.length, iterationNum, "Pomodoro")
		resting = true
		printTransition(writer, fmt.Sprintf("Pomodoro %d finished! Starting rest period",
			iterationNum))

		// Rest Period
		doIteration(writer, opts.restLength, iterationNum, "Rest Period")

		if iterationNum != 0 && iterationNum == opts.iterations {
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
