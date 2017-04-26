package main

import (
	"errors"
	"log"
	"os"

	"github.com/hoffa2/worm/organizer"
	"github.com/hoffa2/worm/segment"
	"github.com/hoffa2/worm/wormgate"
	"github.com/urfave/cli"
)

var (
	ErrportNotSet = errors.New("Port is not set")
)

func main() {
	app := cli.NewApp()
	app.Name = "Awesome worm"
	app.Usage = "Run one of the components"

	app.Commands = []cli.Command{
		{
			Name:  "viz",
			Usage: "run visualizer",
			Action: func(c *cli.Context) error {
				if !c.IsSet("wormport") {
					return errors.New("Wormport flag must be set")
				}
				if !c.IsSet("segmentport") {
					return errors.New("segmentport flag must be set")
				}
				return organizer.Run(c)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "wormport, wp",
					Usage: "Wormagte port (prefix with colon)",
				},
				cli.StringFlag{
					Name:  "segmentport, sp",
					Usage: "segment port (prefix with colon)",
				},
			},
		},
		{
			Name:  "segment",
			Usage: "run segment",
			Action: func(c *cli.Context) error {
				if !c.IsSet("mode") {
					return errors.New("Wormport flag must be set")
				}
				return segment.Run(c)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "wormport, wp",
					Usage: "Wormagte port (prefix with colon)",
				},
				cli.StringFlag{
					Name:  "segmentport, sp",
					Usage: "segment port (prefix with colon)",
				},
				cli.StringFlag{
					Name:  "mode, m",
					Usage: "Spread or Start",
				},
				cli.IntFlag{
					Name:  "target, t",
					Usage: "Inital number of targets (Set only if the segments is the first in the network)",
				},
			},
		},
		{
			Name:  "wormgate",
			Usage: "Starts the wormgate",
			Action: func(c *cli.Context) error {
				return wormgate.Run(c)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "wormport, wp",
					Usage: "Wormagte port (prefix with colon)",
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
