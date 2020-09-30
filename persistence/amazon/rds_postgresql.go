package amazon

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/url"
	"sync"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var awsCreds *credentials.Credentials
var loadCredential sync.Once

//
// hint : https://github.com/califlower/golang-aws-rds-iam-postgres
//
type awsRdsDb struct {
	dbEndpoint string
	dbUser     string
	dbName     string
	dbPassword string
	awsRegion  string
}

func (a *awsRdsDb) Connect(ctx context.Context) (driver.Conn, error) {

	//awsCreds = credentials.NewEnvCredentials()
	// authToken, err := rdsutils.BuildAuthToken(a.dbEndpoint, a.awsRegion, a.dbUser, awsCreds)
	// if err != nil {
	// 	return nil, err
	// }

	psqlURL, err := url.Parse("postgres://")
	if err != nil {
		return nil, err
	}

	psqlURL.Host = a.dbEndpoint
	psqlURL.User = url.UserPassword(a.dbUser, a.dbPassword)
	psqlURL.Path = a.dbName
	q := psqlURL.Query()
	q.Add("sslmode", "disable")

	psqlURL.RawQuery = q.Encode()
	pgxDriver := &stdlib.Driver{}
	connector, err := pgxDriver.OpenConnector(psqlURL.String())
	if err != nil {
		return nil, err
	}

	return connector.Connect(ctx)
}

func (a *awsRdsDb) Driver() driver.Driver {
	return a
}

var DriverNotSupportedErr = errors.New("driver open method not supported")

// driver.Driver interface
func (config *awsRdsDb) Open(name string) (driver.Conn, error) {
	return nil, DriverNotSupportedErr
}

func ConnectRds(dbEndpoint, awsRegion, dbUser, dbName, dbPassword string) *sqlx.DB {
	aDb := &awsRdsDb{
		awsRegion:  awsRegion,
		dbEndpoint: dbEndpoint,
		dbUser:     dbUser,
		dbName:     dbName,
		dbPassword: dbPassword,
	}
	db := sql.OpenDB(aDb)
	return sqlx.NewDb(db, "pgx")
}
