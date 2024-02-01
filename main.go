package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/dolthub/driver"
)

func main() {
	if err := app(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func app() error {
	dbName := "main"
	tableName := "contacts"
	commitName := "John Smith"
	commitEmail := "john.smith@example.com"
	dataPath, err := filepath.Abs("data")
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf(
		"file://%v?commitname=%v&commitemail=%v&database=%v&multistatements=true",
		dataPath,
		commitName,
		commitEmail,
		dbName,
	)

	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		return err
	}

	db, err := sql.Open("dolt", dsn)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	err = ensureDatabase(db, dbName)
	if err != nil {
		return fmt.Errorf("ensure database: %w", err)
	}

	isClean, err := isTableClean(db, dbName, tableName)
	if err != nil {
		return err
	}

	if !isClean {
		return fmt.Errorf("table %q is in your local changes", tableName)
	}

	err = ensureTable(db, dbName, tableName)
	if err != nil {
		return fmt.Errorf("ensure table: %w", err)
	}

	// TODO: Commit table creation

	return nil
}

func ensureDatabase(db *sql.DB, name string) error {
	q := fmt.Sprintf("create database if not exists %v", name)
	return runQuery(db, q)
}

func ensureTable(db *sql.DB, dbName string, name string) error {
	q := fmt.Sprintf(
		"use %v; create table if not exists %v (id integer)",
		dbName,
		name,
	)
	return runQuery(db, q)
}

func isTableClean(db *sql.DB, dbName string, name string) (bool, error) {
	q := fmt.Sprintf(
		"use %v; select table_name from dolt_status where table_name = %q",
		dbName,
		name,
	)
	rows, err := db.Query(q)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
	}

	return count == 0, nil
}

func doltStatus(db *sql.DB) error {
	rows, err := db.Query("select * from dolt_status")
	if err != nil {
		return err
	}
	defer rows.Close()

	return printRows(rows)
}

func runQuery(db *sql.DB, query string) error {
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func printRows(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	fmt.Println(strings.Join(cols, "|"))

	for rows.Next() {
		values := make([]interface{}, len(cols))
		var generic = reflect.TypeOf(values).Elem()
		for i := 0; i < len(cols); i++ {
			values[i] = reflect.New(generic).Interface()
		}

		err = rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		result := bytes.NewBuffer(nil)
		for i := 0; i < len(cols); i++ {
			if i != 0 {
				result.WriteString("|")
			}

			var rawValue = *(values[i].(*interface{}))
			switch val := rawValue.(type) {
			case string:
				result.WriteString(val)
			case int:
				result.WriteString(strconv.FormatInt(int64(val), 10))
			case int8:
				result.WriteString(strconv.FormatInt(int64(val), 10))
			case int16:
				result.WriteString(strconv.FormatInt(int64(val), 10))
			case int32:
				result.WriteString(strconv.FormatInt(int64(val), 10))
			case int64:
				result.WriteString(strconv.FormatInt(val, 10))
			case uint:
				result.WriteString(strconv.FormatUint(uint64(val), 10))
			case uint8:
				result.WriteString(strconv.FormatUint(uint64(val), 10))
			case uint16:
				result.WriteString(strconv.FormatUint(uint64(val), 10))
			case uint32:
				result.WriteString(strconv.FormatUint(uint64(val), 10))
			case uint64:
				result.WriteString(strconv.FormatUint(val, 10))
			case float32:
				result.WriteString(strconv.FormatFloat(float64(val), 'f', 2, 64))
			case float64:
				result.WriteString(strconv.FormatFloat(val, 'f', 2, 64))
			case bool:
				if val {
					result.WriteString("true")
				} else {
					result.WriteString("false")
				}
			case []byte:
				enc := base64.NewEncoder(base64.URLEncoding, result)
				_, err := enc.Write(val)
				return fmt.Errorf("failed to base64 encode blob: %w", err)
			case time.Time:
				timeStr := val.Format(time.RFC3339)
				result.WriteString(timeStr)
			}
		}

		fmt.Println(result.String())
	}

	return nil
}
