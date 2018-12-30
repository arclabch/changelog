// Copyright © 2018 ARClab, Lionel Riem - https://arclab.ch/
//
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

const (
	// Message if no result
	showNoResult = "No entry found.\n"

	// Length needed for the date and the spaces
	dateSpacing = 20

	// Default number of entries to display
	defaultLimit = 5
)

// show fetches and prints entries from the database.
// It returns an error if something goes wrong.
func show(c *cli.Context) error {

	var result []entry
	var count int
	var length int

	desiredWidth := 0
	limit := defaultLimit
	userName := ""

	// Does the DB exist? If not, tell user so
	if !checkDbExists() {
		fmt.Printf(showNoResult)
		return nil
	}

	// Initialise the DB
	err := initDb()
	if err != nil {
		return err
	}

	// Get the number of entries to display, if specified
	if flagLimit := c.Int("last"); flagLimit > 0 {
		limit = flagLimit
	}

	// Check the number of entries to display
	if limit < 1 {
		return errors.New("can't display 0 entry")
	}

	// Check if a username was supplied
	if flagUser := c.String("user"); flagUser != "" {
		userName = flagUser
	}

	if flagWidth := c.Int("width"); flagWidth > 0 {
		desiredWidth = flagWidth
	}

	// Fetch the entries
	result, count, length, err = getEntries(limit, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf(showNoResult)
			return nil
		}
		return err
	}

	// If not result, say so
	if count == 0 {
		fmt.Printf(showNoResult)
		return nil
	}

	// Process and print the entries
	for _, e := range result {
		fmt.Printf("%s  %s  %s\n",
			e.timestamp.Format("2006-01-02T15:04"),
			pad(e.user, length, false), formatMessage(e.message, length+dateSpacing, desiredWidth))
	}

	return nil
}

// pad pads a string with spaces at the left or right of the string
// It returns a string.
func pad(s string, l int, left bool) string {
	length := len([]rune(s))
	diff := l - length
	padding := ""

	if diff > 0 {
		for i := 0; i < diff; i++ {
			padding = padding + " "
		}
	}

	if left == true {
		return padding + s
	}
	return s + padding
}

// formatMessage formats the message for printing on screen.
// It returns a string.
func formatMessage(message string, minus int, width int) string {
	newMessage := ""
	lineLength := 0
	padding := ""

	// Prepare the padding on the left
	for i := 0; i < minus; i++ {
		padding = padding + " "
	}

	// If a width value superior to 0 is supplied, use it. Otherwise, try to get the
	// real terminal size. Substract the minus size from it.
	if width > 0 {
		lineLength = width - minus
	} else {
		lineLength = getTerminalWidth() - minus
	}

	// If the resulting line length is inferior to 10, set it to 10.
	if lineLength < 10 {
		lineLength = 10
	}

	// If the message is shorter than lineLength and doesn't have \n in it
	if len([]rune(message)) <= lineLength && strings.Index(message, "\n") == -1 {
		return message
	}

	// Split the message at \n into slices
	msgSlice := strings.Split(message, "\n")

	// Go through the slices
	for i, s := range msgSlice {
		if i == 0 && len([]rune(s)) <= lineLength {
			newMessage = newMessage + s + "\n"
		} else if len([]rune(s)) <= lineLength {
			newMessage = newMessage + padding + s + "\n"
		} else {
			for j, str := range chunkByWords(s, lineLength) {
				if i == 0 && j == 0 {
					newMessage = newMessage + str + "\n"
				} else {
					newMessage = newMessage + padding + str + "\n"
				}
			}
		}
	}

	// Remove last newline and return the message
	return strings.TrimSuffix(newMessage, "\n")
}

// chunkByWords takes a string and returns a slice of strings cut to maxLength
func chunkByWords(s string, maxLength int) []string {
	sub := ""
	subs := []string{}
	runes := []rune(s)
	l := len(runes)
	index := 0

	// Loop until we reach the end of the index
	for index < l {
		var toIndex int
		toIndex = index + maxLength + 1

		// If target toIndex is still below length of string, process
		if l > toIndex {
			lastSpace := 0

			// Prepare temporary substring
			subRunes := runes[index:toIndex]

			// Search last space in substring
			for pos, r := range subRunes {
				if string(r) == " " {
					lastSpace = pos
				}
			}

			// Use it to extract real substring and reset toIndex
			toIndex = index + lastSpace
			sub = string(runes[index:toIndex])
			subs = append(subs, sub)

		} else {

			// If target toIndex is higher than length, give the remainder of the string
			toIndex = l
			sub = string(runes[index:toIndex])
			subs = append(subs, sub)
		}

		// Prepare next index to last toIndex + 1 character (removes space on linebreak)
		index = toIndex + 1
	}

	return subs
}
