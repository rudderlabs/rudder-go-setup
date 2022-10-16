package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/rudderlabs/rudder-go-setup/setup"
	"github.com/urfave/cli/v2"
)

var App = &cli.App{
	Name:  "go-setup",
	Usage: "setup tools and configs for go project",
	Commands: []*cli.Command{
		Init(),
	},
}

func init() {
	sort.Sort(cli.FlagsByName(App.Flags))
	sort.Sort(cli.CommandsByName(App.Commands))
}

func main() {
	if err := App.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Init() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "init project",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			p := &setup.Project{}
			if err := p.Detect(); err != nil {
				return err
			}

			fmt.Printf("%+v\n", p)
			if err := p.Init(); err != nil {
				return err
			}

			return nil
		},
	}
}
