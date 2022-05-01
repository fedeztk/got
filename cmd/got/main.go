package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fedeztk/got/internal/config"
	"github.com/fedeztk/got/internal/model"
	"github.com/fedeztk/got/pkg/translator"
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
	oneShot := flag.Bool(
		"o",
		false,
		"one shot",
	)
	from := flag.String(
		"f",
		"",
		"from",
	)
	to := flag.String(
		"t",
		"",
		"to",
	)
	flag.Parse()

	switch {
	case *showVersion:
		fmt.Println(gotVersion)

	case *oneShot:
		if *from == "" || *to == "" {
			fmt.Println("from and to are required in one shot mode")
			os.Exit(1)
		}
		response, err := translator.Translate(strings.Join(flag.Args(), " "), *from, *to)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(response.PrettyPrint())

	default:
		model.Run(config.NewConfig())
	}

	os.Exit(0)
}
