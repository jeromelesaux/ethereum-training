package amazon

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/url"
	"sync"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var awsCreds *credentials.Credentials
var loadCredential sync.Once

type awsRdsDb struct {
	dbEndpoint string
	dbUser     string
	dbName     string
	awsRegion  string
}

func (a *awsRdsDb) Connect(ctx context.Context) (driver.Conn, error) {

	awsCreds = credentials.NewEnvCredentials()
	authToken, err := rdsutils.BuildAuthToken(a.dbEndpoint, a.awsRegion, a.dbUser, awsCreds)
	if err != nil {
		return nil, err
	}

	psqlURL, err := url.Parse("postgres://")
	if err != nil {
		return nil, err
	}

	psqlURL.Host = a.dbEndpoint
	psqlURL.User = url.UserPassword(a.dbUser, authToken)
	psqlURL.Path = a.dbName
	q := psqlURL.Query()
	q.Add("sslmode", "true")

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

func ConnectRds(dbEndpoint, awsRegion, dbUser, dbName string) *sqlx.DB {
	aDb := &awsRdsDb{
		awsRegion:  awsRegion,
		dbEndpoint: dbEndpoint,
		dbUser:     dbUser,
		dbName:     dbName,
	}
	db := sql.OpenDB(aDb)
	return sqlx.NewDb(db, "pgx")
}
