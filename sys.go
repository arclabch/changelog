// Copyright © 2018 ARClab, Lionel Riem - https://arclab.ch/
//
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"os"
	"os/user"

	"github.com/olekukonko/ts"
)

// Minimal terminal width if it can't be found or is smaller than that
const minWidth = 40

// checkDbExists checks if the database file exists.
// It returns true or false.
func checkDbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// setDbRights sets the access rights (644) to the database file.
// It returns an error if something goes wrong.
func setDbRights() error {
	return os.Chmod(dbFile, 0644)
}

// userIsRoot checks if the user is root.
// It returns true or false.
func userIsRoot() bool {
	if os.Geteuid() == 0 {
		return true
	}
	return false
}

// getRealUser gets the real login of the user, if possible.
// It returns a string.
func getRealUser() string {

	// Try first with SUDO_USER if user is root
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" && userIsRoot() {
		return sudoUser
	}

	// Check if LOGNAME environment variable exists
	if logname := os.Getenv("LOGNAME"); logname != "" {
		return logname
	}

	// Try to get it from USER
	if user := os.Getenv("USER"); user != "" {
		return user
	}

	// Finally, try to get it from the UID
	if user, err := user.Current(); err == nil {
		return user.Username
	}

	// If all else fails, return an empty string
	return ""
}

// getTerminalWidth gets the width of the terminal.
// Returns an int.
func getTerminalWidth() int {

	// Get the terminal size. In case of error, return default width
	size, err := ts.GetSize()
	if err != nil {
		return minWidth
	}

	// If real size is smaller than default width, return default width
	width := size.Col()
	if width < 40 {
		return minWidth
	}

	// Return terminal width
	return width
}
