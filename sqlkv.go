package sqlkv

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type SqlKv struct {
	db                *sql.DB
	tableName         string
	createTableCalled bool
}

type SqlKvRow struct {
	name  string
	value string
}

func New(db *sql.DB, tableName string) *SqlKv {
	output := new(SqlKv)
	output.db = db
	output.tableName = tableName
	output.createTableCalled = false

	err := output.createTable()
	if err != nil {
		panic(err)
	}

	return output
}

func (this *SqlKv) createTable() error {
	if this.createTableCalled {
		return nil
	}

	this.createTableCalled = true
	_, err := this.db.Exec("CREATE TABLE IF NOT EXISTS " + this.tableName + " (name TEXT NOT NULL PRIMARY KEY, value TEXT)")
	if err != nil {
		return err
	}

	_, err = this.db.Exec("CREATE INDEX name_index ON " + this.tableName + " (name)")
	return err
}

func (this *SqlKv) rowByName(name string) (*SqlKvRow, error) {
	row := new(SqlKvRow)
	query := "SELECT name, `value` FROM " + this.tableName + " WHERE name = ?"
	err := this.db.QueryRow(query, name).Scan(&row.name, &row.value)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return row, nil
}

func (this *SqlKv) String(name string) string {
	row, err := this.rowByName(name)
	if err == nil && row == nil {
		return ""
	}
	if err != nil {
		panic(err)
	}
	return row.value
}

func (this *SqlKv) SetString(name string, value string) {
	row, err := this.rowByName(name)
	var query string

	if row == nil && err == nil {
		query = "INSERT INTO " + this.tableName + " (value, name) VALUES(?, ?)"
	} else {
		query = "UPDATE " + this.tableName + " SET value = ? WHERE name = ?"
	}

	_, err = this.db.Exec(query, value, name)

	if err != nil {
		panic(err)
	}
}

func (this *SqlKv) Int(name string) int {
	s := this.String(name)
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return i
}

func (this *SqlKv) SetInt(name string, value int) {
	s := strconv.Itoa(value)
	this.SetString(name, s)
}

func (this *SqlKv) Float(name string) float32 {
	s := this.String(name)
	if s == "" {
		return 0
	}

	o, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	return float32(o)
}

func (this *SqlKv) SetFloat(name string, value float32) {
	s := strconv.FormatFloat(float64(value), 'g', -1, 32)
	this.SetString(name, s)
}

func (this *SqlKv) Bool(name string) bool {
	s := this.String(name)
	return s == "1" || strings.ToLower(s) == "true"
}

func (this *SqlKv) SetBool(name string, value bool) {
	var s string
	if value {
		s = "1"
	} else {
		s = "0"
	}
	this.SetString(name, s)
}

func (this *SqlKv) Time(name string) time.Time {
	s := this.String(name)
	if s == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}

	return t
}

func (this *SqlKv) SetTime(name string, value time.Time) {
	this.SetString(name, value.Format(time.RFC3339Nano))
}

func (this *SqlKv) Del(name string) {
	query := "DELETE FROM " + this.tableName + " WHERE name = ?"
	_, err := this.db.Exec(query, name)

	if err != nil {
		panic(err)
	}
}

func (this *SqlKv) HasKey(name string) bool {
	row, err := this.rowByName(name)
	if row == nil && err == nil {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}
