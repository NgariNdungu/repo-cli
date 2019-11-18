package main

import (
	// "fmt"
	"os"

	"github.com/mkideal/cli"
)

type argT struct {
	Name string `cli:"name" usage:"omae dare?"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("Hello, %s!\n", argv.Name)
		return nil
	}))
}
