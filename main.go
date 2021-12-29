package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/fedeztk/got/pkg/model"
)

const (
	// help
	translateShell  = "https://github.com/soimort/translate-shell/wiki/Distros"
	transExecutable = "trans"
)

var (
    gotVersion string
)

func main() {
	err := checkDependencies()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	showDoc := flag.Bool(
		"help",
		false,
		"show help",
	)
	showVersion := flag.Bool(
		"version",
		false,
		"show version",
	)
	flag.Parse()
	switch {
	case *showDoc:
		fmt.Println(getHelp())
	case *showVersion:
		fmt.Println(gotVersion)
	default:
		model.Run()
	}
    os.Exit(0)
}

func checkDependencies() error {
	_, err := exec.LookPath(transExecutable)
	if err != nil {
		return fmt.Errorf("please install translate-shell, check: %s", translateShell)
	}
	return nil
}

func getHelp() string {
    return `
        Usage: got [options]

        Options:
            -help, -h                 Show this help message
            -version, -v              Show version
    `
}
