package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "input",
				TakesFile: true,
				Aliases:   []string{"i"},
				Usage:     "path to book yaml file",
				Required:  true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "markdown",
				Usage:   "assemble book in markdown format",
				Aliases: []string{"md", "m"},
				Action:  markdown,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "outdir",
						TakesFile: true,
						Aliases:   []string{"o"},
						Usage:     "output directory",
						Required:  true,
					},
					&cli.BoolFlag{
						Name:    "wordcount",
						Aliases: []string{"w"},
						Usage:   "print word count for each scene",
					},
				},
			},
		},
		Usage: "assemble a book",
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
