package main

import (
	"fmt"
	"os"
	"io"

	"github.com/ngarindungu/repo-cli/repo"
	"github.com/mkideal/cli"
)

// commands to manage repos
// create
// list
// delete
func main() {
	// TODO: register commands
	commands := cli.Root(root, cli.Tree(create), cli.Tree(list), cli.Tree(del))
	if err := commands.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type rootT struct {
	cli.Helper
	Token string `cli:"t,token" usage:"Github Oauth token to use" dft:"$GHUB_TOKEN"`
}

var root = &cli.Command{
	Global: true,
	Argv: func() interface{} { return new(rootT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*rootT)
		ctx.String("Passed in token: %s", argv.Token)
		return nil
	},
}

type createT struct {
	cli.Helper
	Name    string `cli:"n,name" usage:"The name of the repo to work on"`
	Private bool   `cli:"private" usage:"Create a private repo" dft:"false"`
}

var create = &cli.Command{
	Name:    "create",
	Aliases: []string{"new"},
	// Text: "Create a new repository",
	Desc: "Create a new repository",
	Argv: func() interface{} { return new(createT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*createT)
		ctx.String("Creating repo: %s", argv.Name)
		return nil
	},
}

type listT struct {
	// cli.Helper
	First int    `cli:"first" usage:"List the first n repos" dft:"10"`
	Last  int    `cli:"last" usage:"List the last n repos" dft:"10"`
	Order string `cli:"order" usage:"Order in which to fetch repos" dft:"creation"`
}

var list = &cli.Command{
	Name: "list",
	Desc: "Retrieve a user's repositories",
	Argv: func() interface{} { return new(listT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*listT)
		rootArgv := ctx.RootArgv().(*rootT)
		repos := repo.List(rootArgv.Token)
		for _,r := range repos {
			io.WriteString(os.Stdout, r)
		}
		ctx.String("Reading repos ordered by %s", argv.Order)
		return nil
	},
}

type delT struct {
	cli.Helper
	Name string `cli:"n,name" usage:"Name of the repo to delete"`
}

var del = &cli.Command{
	Name:    "delete",
	Aliases: []string{"d", "rm"},
	Desc:    "Delete a repository",
	Argv:    func() interface{} { return new(delT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*delT)
		ctx.String("Deleting repo: %s", argv.Name)
		return nil
	},
}
