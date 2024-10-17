# Simple Redis Mock CLI

This project is a simplified, in-memory Redis-like key-value store mocked as a CLI program in Go. It supports basic Redis commands such as `SET`, `GET`, `DEL`, and others, with optional arguments for certain commands. This is intended to serve as a template for various assignemnts.

This is largely inspired by existing [Redis](https://redis.io/) behavior - please see their docs for a more detailed overview: <https://redis.io/docs/latest/commands>

## Available Commands

### 1. `SET key value [NX | XX] [GET]`

- Sets a string value for the given key. If key already holds a value, it will be overwritten by default. If a non-string value is stored here, an error will be returned.
- Optional flags:
  - `NX`: Only set the key if it does not already exist.
  - `XX`: Only set the key if it already exists.
  - `GET`: Return the old value stored at the key, or `(nil)` if there was no previous value.
- Link: <https://redis.io/docs/latest/commands/set/>
- **Example**:

  ```sh
  redis> SET key1 "hello" NX GET
  ```

### 2. `GET key`

- Retrieves the value associated with the given key, or `(nil)` if it does not exist. If a non-string value is stored here, an error will be returned.
- Link: <https://redis.io/docs/latest/commands/get/>
- **Example**:

  ```sh
  redis> GET key1
  ```

### 3. `DEL key [key...]`

- Deletes the specified keys, ignoring non-existing keys.
- Link: <https://redis.io/docs/latest/commands/del/>
- **Example**:

  ```sh
  redis> DEL key1 key2
  ```

## How to Execute

To run locally, just run

```sh
go run main.go store.go utils.go
```

Alternatively, you can build the project and run the executable:

```sh
make start

# OR

go build -o ./bin/redis-go-cli 
./bin/redis-go-cli
```

Then start entering Redis commands!
