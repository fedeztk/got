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
		"one shot mode, requires -f and -t",
	)
	from := flag.String(
		"f",
		"",
		"language to translate from",
	)
	to := flag.String(
		"t",
		"",
		"language to translate to",
	)
	engine := flag.String(
		"e",
		"",
		"engine, could be: google (default), iciba, reverso, libre\nDeepl is not supported yet and defaults to google",
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
		if *engine == "" {
			*engine = "google"
		}
		response, err := translator.Translate(strings.Join(flag.Args(), " "), *from, *to, *engine)
		if err != nil {
			fmt.Println(model.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}
		fmt.Println(response.PrettyPrint())

	default:
		conf := config.NewConfig()
		if *engine != "" {
			conf.SetEngine(*engine)
		}
		model.Run(conf)
	}

	os.Exit(0)
}
