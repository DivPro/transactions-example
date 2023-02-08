package main

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var err error

	db, err := sql.Open("pgx", "postgres://pg-user:pg-pass@127.0.0.1:5436/pg-db?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("testdata/fixtures"),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	if err := fixtures.Load(); err != nil {
		log.Fatalln(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSample(t *testing.T) {
	t.Log("started")
}
