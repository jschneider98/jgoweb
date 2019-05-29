// +build integration

package db

import (
	"testing"
)

//
func TestDb(t *testing.T) {
	db, err := NewDb()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if db == nil {
		t.Errorf("\nERROR: nil db")
	}

	_, err = db.GetSessionByName("bad_session")

	if err == nil {
		t.Errorf("Invalid Db connection should return error")
	}

	sess, err := db.GetSessionByName("uxt_0000")

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sess == nil {
		t.Errorf("\nERROR: nil db session")
	}

	conn, err := db.GetRandomConn()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if conn == nil {
		t.Errorf("\nERROR: nil db conn")
	}
}
