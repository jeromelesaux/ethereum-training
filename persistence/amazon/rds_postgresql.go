package amazon

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
)

var awsCreds *credentials.Credentials
var loadCredential sync.Once

func ConnectRds(dbEndpoint, awsRegion, dbUser, dbName string, db *sql.DB) (*sql.DB, error) {
	var rdsErr error
	loadCredential.Do(
		func() {
			awsCreds = credentials.NewEnvCredentials()
			authToken, err := rdsutils.BuildAuthToken(dbEndpoint, awsRegion, dbUser, awsCreds)
			if err != nil {
				rdsErr = err
				return
			}
			dnsStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true",
				dbUser, authToken, dbEndpoint, dbName,
			)
			// Use db to perform SQL operations on database
			db, err = sql.Open("postgres", dnsStr)
			if err != nil {
				rdsErr = err
			}
			return
		})

	return db, rdsErr
}
