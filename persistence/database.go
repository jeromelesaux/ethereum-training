package persistence

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jeromelesaux/ethereum-training/persistence/amazon"
	"github.com/jeromelesaux/ethereum-training/persistence/local"
	"github.com/jmoiron/sqlx"
)

var db *sql.DB
var dbx *sqlx.DB
var (
	_dbEndpoint string
	_awsRegion  string
	_dbUser     string
	_dbName     string
	_dbPassword string

	_useSqlite bool
)

func Initialise(useSqlite bool, dbEndpoint, awsRegion, dbUser, dbName, dbpassword string) error {
	_dbEndpoint = dbEndpoint
	_awsRegion = awsRegion
	_dbUser = dbUser
	_dbName = dbName
	_dbPassword = dbpassword
	_useSqlite = useSqlite
	if err := connect(); err != nil {
		return err
	}
	if err := createSchema(); err != nil {
		return err
	}
	return nil
}

func connect() error {
	var err error
	if _useSqlite {
		db, err = local.ConnectSqlite()
		if err != nil {
			return err
		}
	} else {
		dbx = amazon.ConnectRds(_dbEndpoint, _awsRegion, _dbUser, _dbName, _dbPassword)

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

func createSchema() error {
	schema := "create table if not exists 'documents' (" +
		"uid integer primary key autoincrement, " +
		"userid varchar(64), " +
		"created date null, " +
		"document varchar(64), " +
		"checksum varchar(64), " +
		"txhash varchar(64));"
	if !_useSqlite {
		schema = "create table if not exists 'documents' (" +
			"uid serial primary key , " +
			"userid varchar(64), " +
			"created date null, " +
			"document varchar(64), " +
			"checksum varchar(64), " +
			"txhash varchar(64));"
	}
	fmt.Fprintf(os.Stdout, "%s\n", schema)

	if err := connect(); err != nil {
		return err
	}
	var stmt *sql.Stmt
	var err error
	if _useSqlite {
		stmt, err = db.Prepare(schema)
	} else {
		stmt, err = dbx.Prepare(schema)
	}

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
	var err error
	var tx *sql.Tx

	if _useSqlite {
		tx, err = db.Begin()

		if err != nil {
			return err
		}
	} else {
		tx, err = dbx.Begin()
	}
	query := "insert into documents(userid,created,document,checksum,txhash) values('" +
		document.UserID + "',?,'" +
		//	document.Created.String() + ",'" +
		document.DocumentName + "','" +
		document.Checksum + "','" +
		document.TxHash + "'" +
		");"
	fmt.Fprintf(os.Stdout, "%s\n", query)

	insert, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	_, err = insert.Exec(document.Created)
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

	var res *sql.Rows
	if _useSqlite {
		res, err = db.Query(query)
		if err != nil {
			return docs, err
		}
	} else {
		res, err = dbx.Query(query)
		if err != nil {
			return docs, err
		}
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

	var res *sql.Rows
	if _useSqlite {
		res, err = db.Query(query)
		if err != nil {
			return docs, err
		}
	} else {
		res, err = dbx.Query(query)
		if err != nil {
			return docs, err
		}
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

	var res *sql.Rows
	if _useSqlite {
		res, err = db.Query(query)
		if err != nil {
			return docs, err
		}
	} else {
		res, err = dbx.Query(query)
		if err != nil {
			return docs, err
		}
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
