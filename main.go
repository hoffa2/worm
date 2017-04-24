package main

import (
	"errors"
	"log"
	"os"

	"github.com/hoffa2/worm/segment"
	"github.com/hoffa2/worm/visualize"
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
				return visualize.Run(c)
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
				mode := c.String("mode")
				if mode == "spread" {
					return segment.SendSegment(c)
				} else if mode == "start" {
					return segment.StartSegmentServer(c)
				}
				return errors.New("Mode must be either \"mode\" or \"spread\"")
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
