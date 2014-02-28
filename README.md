# An SQL-based key/value store for Golang

[![Build Status](https://travis-ci.org/laurent22/go-sqlkv.png)](https://travis-ci.org/laurent22/go-sqlkv)

SqlKv provides an SQL-based key/value store for Golang. It can work with any of the database types supported by the built-in `database/sql` package.

It can be used, for example, to easily store configuration values in Sqlite, or to build a simple caching system when something more advanced like memcached or Redis is not available.

# Installation

	go get github.com/laurent22/go-sqlkv
	
# Usage

The first step is to initialize a new database connection. The package expects the connection to remain open while being used. For example, using Sqlite:

```go
db, err := sql.Open("sqlite3", "example.db")
if err != nil {
	panic(err)
}
defer db.Close()
```
	
Then create a new SqlKv object and pass it the db connection and the table name. The table will automatically be created if it does not already exist.

```go
store := sqlkv.New(db, "kvstore")
```
	
The value can then be retrived and set using the provided methods:

- `String(name)` / `SetString(name, value)`
- `Int(name)` / `SetInt(name, value)`
- `Float(name)` / `SetFloat(name, value)`
- `Bool(name)` / `SetBool(name, value)`
- `Time(name)` / `SetTime(name, value)`

In order to keep the API simple, all the errors are handled internally when possible. If an error cannot be handled (eg. cannot read or write to the database), the method will panic.

If a key is missing, each Get method will return Golang's default zero value for this type. The zero values are:

- String: ""
- Int: 0
- Float: 0
- Bool: false
- Time: time.Time{} (Test with `time.IsZero()`)

You can use `HasKey` to check if a key really exists. The method `Del` is also available to delete a key.

# API reference

http://godoc.org/github.com/laurent22/go-sqlkv

# Full example

```go
package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	
	_ "github.com/laurent22/go-sqlkv"
	_ "github.com/mattn/go-sqlite3"	
)

func main() {
	os.Remove("example.db")
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	store := sqlkv.New(db, "kvstore")
	
	store.SetString("username", "John")
	fmt.Println(store.String("username"))
	
	store.SetInt("age", 25)
	fmt.Println(store.Int("age"))
	
	store.SetFloat("pi", 3.14)
	fmt.Println(store.Float("pi"))	

	store.SetTime("today", time.Now())
	fmt.Println(store.Time("today"))	
	
	store.SetBool("enabled", true)
	fmt.Println(store.Bool("enabled"))

	fmt.Println(store.HasKey("username"))
	
	store.Del("username")
	fmt.Println(store.String("username"))

	fmt.Println(store.HasKey("username"))	
}
```

# License

MIT
