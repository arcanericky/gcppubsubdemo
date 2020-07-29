package gcppubsubdemo

import "os"

func ExampleLogLevel() {
	// Enable error and debug logging
	LogLevel(DebugLogLevel)
}

func ExampleLogOutput() {
	writer, _ := os.Create("logfile.log")
	LogOutput(writer)
}
