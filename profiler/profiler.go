package profiler

import (
	"fmt"
	"time"
)

// Profiler : Define the profiler struct
type Profiler struct {
	timeBefore time.Time
	timeAfter  time.Time
	timeTaken  time.Duration
}

var enableProfiling = false
var logProfiling = false

func SetEnableProfiling(value bool) { enableProfiling = value }
func SetLogProfiling(value bool)    { logProfiling = value }

// NewProfiler : Constructor (in Go, we create a function to initialize the struct)
func NewProfiler() *Profiler {
	return &Profiler{
		timeBefore: time.Now(), // Initialize timeBefore with the current time
	}
}

// Reset method to reset the timeBefore
func (p *Profiler) Reset() {
	p.timeBefore = time.Now()
}

// EndRecord method to calculate and return the time taken
func (p *Profiler) EndRecord(message string) string {
	if enableProfiling {
		p.timeAfter = time.Now()                    // Capture the current time
		p.timeTaken = p.timeAfter.Sub(p.timeBefore) // Calculate the time difference

		result := fmt.Sprintf("%s : %v", message, p.timeTaken)

		// Optionally log the result
		if logProfiling {
			fmt.Println(result)
		}

		return fmt.Sprintf("%v", p.timeTaken)
	}

	return ""
}
