// Copyright © 2018 ARClab, Lionel Riem - https://arclab.ch/
//
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var (

	// Where is the database?
	dbFile = "/var/log/changelog.db"

	// Leave those variables as is, will be replaced when built
	version = "unknown"
	kernel  = "unknown"
	machine = "unknown"
	built   = "unknown"
)

func main() {

	// Check if the version is properly set
	if version == "unknown" {
		log.Fatal(errors.New("please follow the instructions in README to build this software"))
		os.Exit(1)
	}
	// Prepare the CLI app
	app := cli.NewApp()
	app.Name = "changelog"
	app.Usage = "Keep a journal of maintenance, upgrades and changes done to a system."
	app.Version = version + " (" + kernel + " " + machine + ", built " + built + ")"

	// Add the commands
	app.Commands = []cli.Command{

		// "Add" command definition
		{
			Name:      "add",
			Usage:     "Add a new entry to the changelog (root only)",
			ArgsUsage: "[\"log entry\"]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "stdin, s",
					Usage: "Take the log entry from stdin instead of from argument",
				},
				cli.StringFlag{
					Name:  "user, u",
					Usage: "Set the entry to user `USER`",
				},
			},
			Action: func(c *cli.Context) error {
				return add(c)
			},
		},

		// "Show" command definition
		{
			Name:      "show",
			Usage:     "Show entries from the changelog",
			ArgsUsage: " ", // Empty for no arguments in help
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "last, l",
					Usage: "Show the last `X` entries",
				},
				cli.StringFlag{
					Name:  "user, u",
					Usage: "Show the entries of user `USER`",
				},
				cli.StringFlag{
					Name:  "width, w",
					Usage: "Format the output to a width of `X` characters",
				},
			},
			Action: func(c *cli.Context) error {
				return show(c)
			},
		},
	}

	// Manage not found commands
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Printf("No matching command '%s'\n\n", command)
		cli.ShowAppHelp(c)
	}

	// Start the app
	err := app.Run(os.Args)

	// If the DB is open, close it.
	if dbOpen {
		db.Close()
	}

	// Goodbye.
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}
