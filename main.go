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
- LPUSH key value [value...]    : Add one or more values to the beginning of the list
- LPOP key                      : Pop the first element from a list
- LLEN key                      : Get the length of a list at key
- LRANGE key start stop         : Get elements of the list stored at key from start to stop
- HSET key field value          : Set a field in a hash
- HGET key field                : Get a field value from a hash
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

		case "LPUSH":
			if len(args) < 3 {
				fmt.Println("(error) wrong number of arguments for 'lpush' command")
				continue
			}

			key, values := args[1], args[2:]

			if count, err := store.LPush(key, values...); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("(integer)", count)
			}

		case "LPOP":
			if len(args) != 2 {
				fmt.Println("(error) wrong number of arguments for 'lpop' command")
				continue
			}

			key := args[1]

			// TODO: add count support
			value, found, err := store.LPop(key)
			if err != nil {
				fmt.Println(err)
			} else if !found {
				fmt.Println("(nil)")
			} else {
				fmt.Println(value)
			}

		case "LLEN":
			if len(args) != 2 {
				fmt.Println("(error) wrong number of arguments for 'llen' command")
				continue
			}

			key := args[1]
			length, err := store.LLen(key)
			if err != nil {
				fmt.Println("(error)", err)
			} else {
				fmt.Println("(integer)", length)
			}

		case "LRANGE":
			if len(args) != 4 {
				fmt.Println("(error) wrong number of arguments for 'lrange' command")
				continue
			}

			key := args[1]
			start, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("(error) ERR value is not an integer or out of range")
				continue
			}

			stop, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println("(error) ERR value is not an integer or out of range")
				continue
			}

			result, err := store.LRange(key, start, stop)
			if err != nil {
				fmt.Println(err)
			} else if len(result) == 0 {
				fmt.Println("(empty array)")
			} else {
				for i, value := range result {
					fmt.Printf("%d) \"%s\"\n", i+1, value)
				}
			}

		case "HSET":
			if len(args) != 4 {
				fmt.Println("(error) wrong number of arguments for 'hset' command")
				continue
			}

			key, field, value := args[1], args[2], args[3]

			if err := store.HSet(key, field, value); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("(integer) 1")
			}

		case "HGET":
			if len(args) != 3 {
				fmt.Println("(error) wrong number of arguments for 'hget' command")
				continue
			}

			key, field := args[1], args[2]

			value, found, err := store.HGet(key, field)
			if err != nil {
				fmt.Println(err)
			} else if !found {
				fmt.Println("(nil)")
			} else {
				fmt.Println(value)
			}

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
