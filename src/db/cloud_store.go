package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	connectionName = mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user           = mustGetenv("CLOUDSQL_USER")
	password       = os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}

// OpenStore creates store from provided sqllite db path
func OpenCloudStore() (*Store, error) {

	// dbURI = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", dbTCPHost, dbUser, dbPwd, dbPort, dbName)
	// https://cloud.google.com/sql/docs/mysql/connect-app-engine-standard
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", user, password, "/cloudsql", connectionName, "card"))

	if err != nil {
		return nil, err
	}
	// todo: check for tables
	return &Store{db}, nil
}

func (s *Store) CloudDBs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	rows, err := s.db.Query("SHOW DATABASES")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not query db: %v", err), 500)
		return
	}
	defer rows.Close()

	buf := bytes.NewBufferString("Databases:\n")
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			http.Error(w, fmt.Sprintf("Could not scan result: %v", err), 500)
			return
		}
		fmt.Fprintf(buf, "- %s\n", dbName)
	}
	w.Write(buf.Bytes())
}
