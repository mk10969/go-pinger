package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

func main() {
	host := flag.String("host", "127.0.0.1", "host")
	port := flag.String("port", "3306", "port")
	user := flag.String("user", "docker", "user")
	pass := flag.String("pass", "docker", "pass")
	name := flag.String("name", "test_database", "database name")
	rawTimeout := flag.String("timeout", "60s", "connection wait timeout")

	flag.Parse()

	cfg := mysql.Config{
		User:                 *user,
		Passwd:               *pass,
		Net:                  "tcp",
		Addr:                 *host + ":" + *port,
		DBName:               *name,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	exitIfErr(err, 1)

	timeout, err := time.ParseDuration(*rawTimeout)
	exitIfErr(err, 1)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = ping(ctx, db)
	if err == nil {
		fmt.Println("ping ok")
	} else {
		fmt.Fprintf(os.Stderr, "\n%s\n", err.Error())
		os.Exit(2)
	}

}

func ping(ctx context.Context, db *sql.DB) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := db.Ping()
		if err == nil {
			return nil
		}
		fmt.Printf(".")
		time.Sleep(time.Second)
	}
}

func exitIfErr(err error, code int) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(code)
	}
}
