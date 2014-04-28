package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	var distanceFlag = flag.Int64("distance", 0, "distance to row in meters")
	var durationFlag = flag.Duration("duration", 0, "duration to row (e.g. 1800s or 45m)")
	var debug = flag.Bool("debug", false, "debug communication data packets")

	flag.Parse()
	if !flag.Parsed() {
		log.Fatal(flag.ErrHelp)
	}

	// client mode
	log.Println("CLI mode - Press CTRL+C to interrupt")

	logCallback := func(ch chan Event) {
		for {
			event := <-ch
			fmt.Printf("%d %s:%d\n", event.time, event.label, event.value)
		}
	}

	workout := NewWorkout(*durationFlag, *distanceFlag)

	s4 := NewS4(logCallback, *debug)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Print("Received signal " + sig.String())
			s4.Exit()
			os.Exit(0)
		}
	}()

	// TODO allow goroutine channel to interrupt workout
	s4.Run(workout)
}
