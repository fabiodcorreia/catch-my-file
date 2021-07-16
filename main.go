package main

import (
	"os"

	"github.com/fabiodcorreia/catch-my-file/cmd/frontend"
)

func main() {
	f := frontend.New()

	if err := f.Run(); err != nil {
		os.Exit(1)
	}

}
