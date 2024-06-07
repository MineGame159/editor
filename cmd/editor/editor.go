package main

import (
	"editor/internal/editor"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "editor",
		Description: "Simple text editor.",
		Action: func(c *cli.Context) error {
			e, err := editor.New()
			if err != nil {
				return err
			}

			err = e.LoadBuffer(c.Args().First())
			if err != nil {
				return err
			}

			return e.Run()
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
