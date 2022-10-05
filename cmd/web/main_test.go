package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	db, err := run()
	if err != nil {
		t.Error("failed run()")
	}
	//fmt.Println("Connected to database")
	defer db.Close()
}
