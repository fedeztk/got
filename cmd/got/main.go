package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fedeztk/got/internal/config"
	"github.com/fedeztk/got/internal/model"
)

var (
	gotVersion string
)

func main() {
	showVersion := flag.Bool(
		"v",
		false,
		"show version",
	)
	flag.Parse()

	switch {
	case *showVersion:
		fmt.Println(gotVersion)
	default:
		model.Run(config.NewConfig())
	}

	os.Exit(0)
}
