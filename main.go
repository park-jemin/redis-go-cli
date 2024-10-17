package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Instructions for the CLI
var instructions = `
Simple Redis CLI. Commands: 
- SET key value [NX | XX] [GET] : Set a key to a string value
- GET key                       : Get the string value of key, or nil if the key does not exist
- DEL key [key...]              : Deletes the specified keys, ignoring non-existing keys
- HELP                          : Show available commands
- QUIT                          : Close the CLI

See the README or the Redis command documentation for more info on behavior and options.
`

func main() {
	store := NewStore()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(instructions)

	defer fmt.Println("Exiting...")

	for {
		fmt.Print("redis> ")
		input, _ := reader.ReadString('\n')
		args := strings.Fields(strings.TrimSpace(input))

		if len(args) == 0 {
			continue
		}

		command := strings.ToUpper(args[0])

		switch command {
		case "SET":
			// TODO: implement TTL options
			if len(args) < 2 || len(args) > 4 {
				fmt.Println("(error) wrong number of arguments for 'set' command")
				continue
			}

			key, value := args[1], args[2]
			options, err := ParseSetOptions(args[3:])
			if err != nil {
				fmt.Println("(error)", err)
				continue
			}

			prev, exists, err := store.Get(key)
			if err != nil {
				fmt.Println("(error)", err)
				continue
			}

			// Lol I should clean this up
			if options.NX && exists {
				if options.GET {
					fmt.Printf("\"%s\"\n", prev)
				} else {
					fmt.Println("(nil)")
				}
				continue
			}

			if options.XX && !exists {
				fmt.Println("(nil)")
				continue
			}

			if err := store.Set(key, value); err != nil {
				fmt.Println("(error)", err)
				continue
			}

			if !options.GET {
				fmt.Println("OK")
			} else if exists {
				fmt.Printf("\"%s\"\n", prev)
			} else {
				fmt.Println("(nil)")
			}

		case "GET":
			if len(args) != 2 {
				fmt.Println("(error) wrong number of arguments for 'get' command")
				continue
			}

			if value, found, err := store.Get(args[1]); err != nil {
				fmt.Println("(error)", err)
			} else if !found {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", value)
			}

		case "DEL":
			if len(args) < 2 {
				fmt.Println("(error) wrong number of arguments for 'del' command")
				continue
			}

			deleted := store.Del(args[1:]...)
			fmt.Println("(integer)", deleted)

		case "HELP":
			fmt.Println(instructions)

		case "EXIT":
			return
		case "QUIT":
			return
		case "Q":
			return

		default:
			fmt.Printf("(error) unknown command '%v'\n", args[0])
		}
	}
}
