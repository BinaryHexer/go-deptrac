package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BinaryHexer/go-deptrac/validator"
)

func main() {
	debug := flag.Bool("debug", true, "debug mode")
	if *debug {
		validator.Log.SetOutput(os.Stderr)
	}

	config := validator.ParseConfig()
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
