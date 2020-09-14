package main

import (
	"os"

	"github.com/dimtass/adcsamplerbeat/cmd"

	_ "github.com/dimtass/adcsamplerbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
