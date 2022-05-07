package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/frioux/shortlinks/shortlinks"
	"github.com/frioux/shortlinks/sqlitestorage"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var publicListen, listen, dsn string

	fs.StringVar(&listen, "listen", ":8080", "address to listen on for read-write server")
	fs.StringVar(&publicListen, "public-listen", "", "address to listen on for public server")
	fs.StringVar(&dsn, "db", "file:db.db", "database file")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	db, err := sqlitestorage.Connect(dsn)
	if err != nil {
		return err
	}

	s := shortlinks.Server{DB: db}

	if publicListen != "" {
		go s.PublicListenAndServe(publicListen)
	}

	return s.ListenAndServe(listen)
}
