# Simple Redis Mock CLI

This project is a simplified, in-memory Redis-like key-value store mocked as a CLI program in Go. It supports basic Redis commands such as `SET`, `GET`, `DEL`, and others, with optional arguments for certain commands. This is intended to serve as a template for various assignemnts.

This is largely inspired by existing [Redis](https://redis.io/) behavior - please see their docs for a more detailed overview: <https://redis.io/docs/latest/commands>

## Usage

To run locally, run

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

### 4. `LPUSH key value [value ...]`

- Pushes one or more values onto the left end of a list stored at  `key`. If a key does not exist, an empty list is created before pushing. If a non-list value is stored here, an error will be returned.
- Link: <https://redis.io/docs/latest/commands/lpush/>
- **Example**:

  ```sh
  redis> LPUSH mylist "one" "two" "three"
  ```

### 5. `LPOP key`

- Get the length of a list at `key`, 0 if the key does not exist. If a non-list value is stored here, an error will be returned.
- Link: <https://redis.io/docs/latest/commands/lpop/>
- **Example**:

  ```sh
  redis> LPOP mylist
  ```

### 6. `LLEN key`

- Removes and returns the first element of the list stored at `key`. If a non-list value is stored here, an error will be returned.
- Link: <https://redis.io/docs/latest/commands/llen/>
- **Example**:

  ```sh
  redis> LLEN mylist
  ```

### 7. `LRANGE key start stop`

- Returns a range of elements from the list stored at `key`, between the `start` and `stop` indexes (inclusive). Supports negative indexes to count from the end of the list. If a non-list value is stored here, an error will be returned.
- Link: <https://redis.io/docs/latest/commands/lrange/>
- **Example**:

  ```sh
  redis> LRANGE mylist 0 -1 # Gets all items in the list from the head
  redis> LRANGE mylist -3 2 # Gets 3 elemnts from the third rightmost element to the third leftmost
  redis> LRANGE mylist -3 -1 # Gets the last 3 elements
  ```

## Testing

To run tests, run

```sh
make test

# OR

go test -v
```
