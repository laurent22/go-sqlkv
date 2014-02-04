# SQL-based key/value store for Golang

SqlKv provides an SQL-based key/value store for Golang. It can work with any of the database type supported by the built-in database/sql package.

It can be used, for example, to easily store configuration values in Sqlite.

It can also be used to build a simple caching system when something like memcached or Redis is not available.

# Installation

	go get github.com/laurent22/go-sqlkv
	
# Usage

The first step is to initialize a new database connection. The package expects the connection to remain open while being used.

	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
Then create a new SqlKv object:

	store := sqlkv.New(db, "kvstore")
	
You can then Get/Set values using the provided methods:

- GetString / SetString
- GetInt / SetInt
- GetFloat / SetFloat
- GetBool / SetBool
- GetTime / SetTime

In order to keep the API simple, all the errors are handled internally when possible. If an error cannot be handled (eg. cannot read or write to the database), the methods will panic.

If a key is missing, each method will return Golang's default zero value for this type. The zero values are:

- String: ""
- Int: 0
- Float: 0
- Bool: false
- Time: time.Time{} (Test with time.IsZero())

# Full example

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
		fmt.Println(store.GetString("username"))
		
		store.SetInt("age", 25)
		fmt.Println(store.GetInt("age"))
		
		store.SetFloat("pi", 3.14)
		fmt.Println(store.GetFloat("pi"))	

		store.SetTime("today", time.Now())
		fmt.Println(store.GetTime("today"))	
		
		store.SetBool("enabled", true)
		fmt.Println(store.GetBool("enabled"))

		fmt.Println(store.HasKey("username"))
		
		store.Del("username")
		fmt.Println(store.GetString("username"))

		fmt.Println(store.HasKey("username"))	
	}

# License

MIT