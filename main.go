package main

import (
	"os"

	"github.com/fabiodcorreia/catch-my-file/pkg/catchmyfile"
	"github.com/fabiodcorreia/catch-my-file/pkg/clog"
)

const port = 8822

func main() {
	/*
		// To monitor the application
		go func() {
			var ms runtime.MemStats
			for {
				runtime.ReadMemStats(&ms)
				clog.Info("Memory: %d - Goroutines: %d", ms.Sys, runtime.NumGoroutine())
				time.Sleep(5 * time.Second)
			}

		}()
	*/

	app := catchmyfile.New(port)

	defer clog.Close()
	clog.Info("========== Catch My File - Started ==========")
	clog.Info("Logging to file: %s", clog.LogFile())

	if err := app.Run(); err != nil {
		clog.Error(err)
		if cErr := clog.Close(); cErr != nil {
			os.Exit(2)
		}
		os.Exit(1)
	}

	clog.Info("========== Catch My File - Finish ===========")
}
