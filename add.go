// Copyright © 2018 ARClab, Lionel Riem - https://arclab.ch/
//
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/urfave/cli"
)

// add adds a new entry to the database.
// It returns an error if something goes wrong.
func add(c *cli.Context) error {
	var user string
	var message string
	var err error

	// Start by checking if we are root
	if !userIsRoot() {
		return errors.New("please add new entries with root privileges")
	}

	// Initialise the DB
	err = initDb()
	if err != nil {
		return err
	}

	// Check if a username was supplied, otherwise get it from system
	if userValue := c.String("user"); userValue != "" {
		user = userValue
	} else {
		user = getRealUser()
	}

	// Check if there is a message in arguments, or try to get it from stdin, or give up
	if c.NArg() > 0 {
		message = ""
		nbArgs := c.NArg()
		for i := 0; i < nbArgs; i++ {
			message = message + c.Args().Get(i)
			if i < nbArgs-1 {
				message = message + " "
			}
		}
	} else if c.Bool("stdin") {
		message, err = readStdin()
		if err != nil {
			return err
		}
	} else {
		return errors.New("please add a message to create an entry")
	}

	// Add the entry
	err = addEntry(user, message)
	if err != nil {
		return err
	}

	return nil
}

// readStdin reads the entry from Stdin.
// It returns the read string and an Error variable.
func readStdin() (string, error) {
	bytes, err := ioutil.ReadAll(os.Stdin)
	msg := strings.TrimSuffix(string(bytes), "\n")
	return msg, err
}
