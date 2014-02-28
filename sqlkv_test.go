package sqlkv

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func getStore() *SqlKv {
	var err error
	var db *sql.DB

	os.Mkdir("test", 0777)

	os.Remove("test/database.db")
	db, err = sql.Open("sqlite3", "test/database.db")
	if err != nil {
		panic(err)
	}

	return New(db, "kvstore")
}

func clearStore(store *SqlKv) {
	store.db.Close()
	os.RemoveAll("test")
}

func panicHandler(t *testing.T, message string) {
	if r := recover(); r == nil {
		t.Errorf("%s: Expected call to panic, but it didn't", message)
	}
}

func noPanicHandler(t *testing.T, message string) {
	if r := recover(); r != nil {
		t.Errorf("%s: Expected call to not panic, but it did", message)
	}
}

func Test_rowByName(t *testing.T) {
	store := getStore()
	defer clearStore(store)
	
	row, err := store.rowByName("test")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	if row != nil {
		t.Error("Expected no data but got", row)
	}
	
	store.SetString("name", "lau")
	row, err = store.rowByName("name")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	if row == nil {
		t.Error("Expected data but got nil")
	}
}

func Test_GetSetString(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	store.SetString("test", "abcd")
	value := store.String("test")
	if value != "abcd" {
		t.Errorf("Expected 'abcd', got '%s'", value)
	}

	store.SetString("test", "1234")
	value = store.String("test")
	if value != "1234" {
		t.Errorf("Expected '1234', got '%s'", value)
	}

	store.db.Close()

	defer panicHandler(t, "String: database is closed")
	store.String("test")

	defer panicHandler(t, "SetString: database is closed")
	store.SetString("test", "panic")
}

func Test_GetSetInt(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	store.SetInt("test", 1234)
	i := store.Int("test")
	if i != 1234 {
		t.Errorf("Expected 1234, got %d", i)
	}
	
	i = store.Int("doesntexist")
	if i != 0 {
		t.Errorf("Expected 0, got %d", i)
	}

	store.SetString("test", "abcd")
	defer panicHandler(t, "Int: not a number")
	store.Int("test")
}

func Test_GetSetFloat(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	store.SetFloat("test", 1234.567)
	f := store.Float("test")
	if f != 1234.567 {
		t.Errorf("Expected 1234.567, got %f", f)
	}
	
	f = store.Float("doesntexist")
	if f != 0 {
		t.Errorf("Expected 0, got %f", f)
	}
	
	store.SetString("test", "abcd")

	defer panicHandler(t, "Float: not a number")
	store.Float("test")
}

func Test_GetSetBool(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	b := store.Bool("nothere")
	if b {
		t.Error("Expected false, got true")
	}

	store.SetBool("test", true)
	if !store.Bool("test") {
		t.Error("Expected true, got false")
	}

	store.SetBool("test", false)
	if store.Bool("test") {
		t.Error("Expected false, got true")
	}
}

func Test_GetSetTime(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	v := store.Time("nothere")
	if !v.IsZero() {
		t.Errorf("Expected zero value, got %s", t)
	}

	now := time.Now()
	store.SetTime("test", now)
	v = store.Time("test")
	if v.Format(time.Stamp) != now.Format(time.Stamp) {
		t.Errorf("Expected %s, got %s", now, v)
	}

	store.SetString("test", "not a date")

	defer panicHandler(t, "Time: not a date")
	store.Time("test")
}

func Test_Del(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	defer noPanicHandler(t, "Del: deleting non-existant key")
	store.Del("blabla")

	store.SetString("test", "abcd")
	store.Del("test")

	value := store.String("test")
	if value != "" {
		t.Errorf("Expected '', got '%s'", value)
	}

	store.db.Close()
	defer panicHandler(t, "Del: database is closed")
	store.Del("test")
}

func Test_HasKey(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	if store.HasKey("test") {
		t.Error("Expected false, got true")
	}

	store.SetString("test", "abcd")
	if !store.HasKey("test") {
		t.Error("Expected true, got false")
	}

	store.SetString("test", "")
	if !store.HasKey("test") {
		t.Error("Expected true, got false")
	}

	store.Del("test")
	if store.HasKey("test") {
		t.Error("Expected false, got true")
	}
}

func Test_All(t *testing.T) {
	store := getStore()
	defer clearStore(store)
	
	store.SetString("test", "abcd")
	store.SetInt("num", 1234)
	store.SetBool("boolean", true)
	
	all := store.All()
	if len(all) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(all))
	}
	
	expected := []SqlKvRow{
		SqlKvRow{ Name: "test", Value: "abcd", },
		SqlKvRow{ Name: "num", Value: "1234", },
		SqlKvRow{ Name: "boolean", Value: "1", },
	}
	
	for _, expectedKv := range expected {
		var found bool
		for _, kv := range all {
			if kv.Name == expectedKv.Name && kv.Value == expectedKv.Value {
				found = true
				break
			}
		}
		if !found {
			t.Error("Not found:", expectedKv)
		}
	}
	
	store.Clear()
	all = store.All()
	if len(all) != 0 {
		t.Errorf("Expected 0 rows")
	}
}

func Test_Placeholder(t *testing.T) {
	store := getStore()
	defer clearStore(store)
	
	if store.placeholder(0) != "?" || store.placeholder(1) != "?" {
		t.Error("Incorrect placeholder")
	}
	
	store.SetDriverName("postgres")
	if store.placeholder(1) != "$1" || store.placeholder(2) != "$2" {
		t.Error("Incorrect placeholder")
	}

	store.SetDriverName("sqlite3")	
	if store.placeholder(0) != "?" || store.placeholder(1) != "?" {
		t.Error("Incorrect placeholder")
	}
}
