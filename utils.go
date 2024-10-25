package main

import (
	"errors"
	"strings"
)

// SetOptions Optional arguments for SET
type SetOptions struct {
	NX  bool // Only set the key if it does not already exist
	XX  bool // Only set the key if it already exists
	GET bool // Return the previous value when setting the key

	// TODO: Implement these options
	// EX      int64 // seconds
	// PX      int64 // milliseconds
	// EXAT    int64 // unix time in seconds
	// PXAT    int64 // unix time in milliseconds
	// KEEPTTL bool
}

// ParseSetOptions parses optional arguments for the SET command
func ParseSetOptions(args []string) (*SetOptions, error) {
	options := &SetOptions{}

	for _, arg := range args {
		switch strings.ToUpper(arg) {
		case "NX":
			options.NX = true
		case "XX":
			options.XX = true
		case "GET":
			options.GET = true
		default:
			return nil, errors.New("ERR syntax error")
		}
	}

	if options.NX && options.XX {
		return nil, errors.New("ERR syntax error")
	}

	return options, nil
}
