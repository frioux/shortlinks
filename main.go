package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/frioux/shortlinks/auth/tailscaleauth"
	"github.com/frioux/shortlinks/shortlinks"
	"github.com/frioux/shortlinks/storage/dynamodbstorage"
	"github.com/frioux/shortlinks/storage/sqlitestorage"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var (
		publicListen, listen, dsn string
		tailscale, dynamodb       bool
	)

	fs.StringVar(&listen, "listen", ":8080", "address to listen on for read-write server")
	fs.StringVar(&publicListen, "public-listen", "", "address to listen on for public server")
	fs.StringVar(&dsn, "db", "file:db.db", "database file")
	fs.BoolVar(&tailscale, "tailscale", false, "enable tailscale auth for read-write server")
	fs.BoolVar(&dynamodb, "dynamodb", false, "enable dynamodb for storage")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	var (
		db  shortlinks.DB
		err error
	)
	if dynamodb {
		db, err = dynamodbstorage.NewClient()
		if err != nil {
			return err
		}
	} else {
		db, err = sqlitestorage.Connect(dsn)
		if err != nil {
			return err
		}
	}

	s := shortlinks.Server{DB: db}
	if tailscale {
		s.Auth = tailscaleauth.Auther{}
	}

	if publicListen != "" {
		go s.PublicListenAndServe(publicListen)
	}

	return s.ListenAndServe(listen)
}
