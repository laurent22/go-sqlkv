package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	
	sqlkv ".."
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