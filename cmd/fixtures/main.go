package main

import (
	"database/sql"
	"flag"
	"github.com/divpro/transactions-example/internal/config"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config.yml", "Configuration file name")
	flag.Parse()
}

func main() {
	textHandler := slog.NewTextHandler(os.Stdout)
	logger := slog.New(textHandler)

	var err error

	f, err := os.Open(configPath)
	if err != nil {
		logger.Error("open configuration file", err, configPath)
		return
	}
	var conf config.Config
	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		logger.Error("parse configuration file", err, configPath)
		return
	}

	db, err := sql.Open("pgx", conf.DB.DSN())
	if err != nil {
		logger.Error("open db", err, conf.DB.DSN())
		return
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
}
