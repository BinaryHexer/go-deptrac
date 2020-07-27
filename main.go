package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BinaryHexer/go-deptrac/validator"
)

func main() {
	debug := flag.Bool("debug", true, "debug mode")
	if *debug {
		validator.Log.SetOutput(os.Stderr)
	}
	flag.Parse()
	var configPath string
	if len(flag.Args()) > 1 {
		configPath = flag.Args()[1]
	} else {
		var err error
		configPath, err = filepath.Abs(flag.Arg(0))
		if err != nil {
			panic(err)
		}
	}

	config := validator.ParseConfig([]string{configPath})
	v := validator.NewValidator(config)

	fmt.Printf("[go-deptrac] checking %s\n", config.Paths[0])
	isValid, errors, err := v.Validate(true)
	if err != nil {
		panic(err)
	}

	if !isValid {
		for _, err := range errors {
			fmt.Println(err.Error())
		}

		os.Exit(1)
	}

	os.Exit(0)
}
