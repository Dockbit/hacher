package main

import (
	"log"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
)

const (
	dirMode os.FileMode = 0755 // Default mode for directories
)

var (
	logger  = log.New(os.Stderr, "", 0)
	verbose = false

	CachePath = os.Getenv("HACHER_PATH")
	CacheKeep = 3
)

func initClient() {
	if len(CachePath) < 1 {
		printFatal("Env variable HACHER_PATH is not set. Point it to the cache directory.")
	}
	if len(os.Getenv("HACHER_KEEP")) > 0 {
		if i, err := strconv.Atoi(os.Getenv("HACHER_KEEP")); err == nil {
			CacheKeep = i
		}
	}
}

func main() {

	initClient()

	app := cli.NewApp()

	app.Name = "Hacher"
	app.Usage = "A simple CLI tool to cache project artifacts."
	app.Author = "The Dockbit Team"
	app.Email = "team@dockbit.com"
	app.Version = "0.1.0"

	// Alphabetically ordered list of commands
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, x",
			Usage: "Verbose mode",
		},
	}

	app.Commands = []cli.Command{

		{
			Name:   "get",
			Usage:  "Gets cache content by a given key.",
			Action: cmdGet,
			Flags: []cli.Flag{

				cli.StringFlag{
					Name:  "k, key",
					Usage: "Cache key",
				},

				cli.StringFlag{
					Name:  "f, file",
					Usage: "Path to comma-separated depdendency file(s) to track for changes.",
				},
			},
		},
		{
			Name:   "set",
			Usage:  "Saves cache content for a given key.",
			Action: cmdSet,
			Flags: []cli.Flag{

				cli.StringFlag{
					Name:  "k, key",
					Usage: "Cache key",
				},

				cli.StringFlag{
					Name:  "f, file",
					Usage: "Path to comma-separated depdendency file(s) to track for changes.",
				},
			},
		},
	}

	app.Run(os.Args)
}
