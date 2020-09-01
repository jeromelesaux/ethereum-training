package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialise() error {
	if err := connect(); err != nil {
		return err
	}
	if err := createSchema(); err != nil {
		return err
	}
	return nil
}

type Document struct {
	UserID       string
	Created      time.Time
	DocumentName string
	Checksum     string
	TxHash       string
}

func NewDocument(userId string, created time.Time, documentName, checksum, txhash string) *Document {
	return &Document{
		UserID:       userId,
		Created:      created,
		DocumentName: documentName,
		Checksum:     checksum,
		TxHash:       txhash,
	}
}

func connect() (err error) {
	db, err = sql.Open("sqlite3", "ethereum.db")
	if err != nil {
		return err
	}
	return nil
}

func createSchema() error {
	schema := "create table if not exists 'documents' (" +
		"uid integer primary key autoincrement, " +
		"userid varchar(64), " +
		"created date null, " +
		"document varchar(64), " +
		"checksum varchar(64), " +
		"txhash varchar(64));"

	fmt.Fprintf(os.Stdout, "%s\n", schema)

	if err := connect(); err != nil {
		return err
	}
	stmt, err := db.Prepare(schema)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return err
}

func InsertDocument(document *Document) error {
	if err := connect(); err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	query := "insert into documents(userid,created,document,checksum,txhash) values('" +
		document.UserID + "','" +
		document.Created.String() + "','" +
		document.DocumentName + "','" +
		document.Checksum + "','" +
		document.TxHash + "'" +
		");"
	fmt.Fprintf(os.Stdout, "%s\n", query)

	insert, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	_, err = insert.Exec()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetDocuments(userID string) (docs []*Document, err error) {
	docs = make([]*Document, 0)
	if err := connect(); err != nil {
		return docs, err
	}
	query := "select userid,created,document,checksum,txhash from documents where userid = '" +
		userID + "';"

	res, err := db.Query(query)
	if err != nil {
		return docs, err
	}

	for res.Next() {
		var userid, document, checksum, txhash string
		var created time.Time
		err = res.Scan(&userid, &created, &document, &checksum, &txhash)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while getting row from sql (%v)\n", err)
		} else {
			docs = append(docs, NewDocument(
				userid,
				created,
				document,
				checksum,
				txhash,
			))
		}

	}
	res.Close()
	return docs, nil
}

func GetDocumentsByName(filename string) (docs []*Document, err error) {
	docs = make([]*Document, 0)
	if err := connect(); err != nil {
		return docs, err
	}
	query := "select userid,created,document,checksum,txhash from documents where document = '" +
		filename + "';"

	res, err := db.Query(query)
	if err != nil {
		return docs, err
	}

	for res.Next() {
		var userid, document, checksum, txhash string
		var created time.Time
		err = res.Scan(&userid, &created, &document, &checksum, &txhash)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while getting row from sql (%v)\n", err)
		} else {
			docs = append(docs, NewDocument(
				userid,
				created,
				document,
				checksum,
				txhash,
			))
		}

	}
	res.Close()
	return docs, nil
}

func GetDocumentsByChecksum(checksum string) (docs []*Document, err error) {
	docs = make([]*Document, 0)
	if err := connect(); err != nil {
		return docs, err
	}
	query := "select userid,created,document,checksum,txhash from documents where checksum = '" +
		checksum + "';"

	res, err := db.Query(query)
	if err != nil {
		return docs, err
	}

	for res.Next() {
		var userid, document, checksum, txhash string
		var created time.Time
		err = res.Scan(&userid, &created, &document, &checksum, &txhash)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while getting row from sql (%v)\n", err)
		} else {
			docs = append(docs, NewDocument(
				userid,
				created,
				document,
				checksum,
				txhash,
			))
		}

	}
	res.Close()
	return docs, nil
}