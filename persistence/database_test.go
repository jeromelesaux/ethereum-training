package persistence

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeromelesaux/ethereum-training/token"
)

func TestConnect(t *testing.T) {
	if err := connect(); err != nil {
		t.Fatalf("Expected no error and gets error (%v)\n", err)
	}
}

func TestCreateSchema(t *testing.T) {
	if err := connect(); err != nil {
		t.Fatalf("Expected no error and gets error (%v)\n", err)
	}
	if err := createSchema(); err != nil {
		t.Fatalf("Expected no error and gets error (%v)\n", err)
	}
}

func TestCreateDocument(t *testing.T) {
	if err := connect(); err != nil {
		t.Fatalf("Expected no error and gets error (%v)\n", err)
	}

	d := NewDocument(
		"test1",
		time.Now(),
		"document1",
		"checksum1",
		"txhash1",
	)

	if err := InsertDocument(d); err != nil {
		t.Fatalf("Expected no error and gets error (%v)\n", err)
	}
}

func TestBulkInsert(t *testing.T) {
	if err := connect(); err != nil {
		t.Fatalf("Expected no error and gets error (%v)\n", err)
	}

	for i := 0; i < 1000; i++ {
		checksum, _ := token.RandToken(64)
		txhash, _ := token.RandToken(64)
		userid := fmt.Sprintf("user%d@mail.com", i)
		document := fmt.Sprintf("document%d.doc", i)

		d := NewDocument(
			userid,
			time.Now(),
			document,
			checksum,
			txhash,
		)

		if err := InsertDocument(d); err != nil {
			t.Fatalf("Expected no error and gets error (%v)\n", err)
		}
	}
}

func TestGetDocument(t *testing.T) {
	if err := connect(); err != nil {
		t.Fatalf("Expected no error and gets error (%v)\n", err)
	}
	docs, err := GetDocuments("test1")
	if err != nil {
		t.Fatalf("Expected no error and gets error (%v)\n", err)
	} else {
		t.Log(docs)
	}
}
