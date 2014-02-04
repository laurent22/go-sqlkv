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