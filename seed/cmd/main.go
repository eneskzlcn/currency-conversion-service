package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/eneskzlcn/currency-conversion-service/config"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/eneskzlcn/currency-conversion-service/seed"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func run() error {
	appConfig := config.New()

	postgresDB, err := postgres.New(appConfig.Db)

	if err != nil {
		fmt.Println("error when initializing database", err.Error())
		return err
	}

	var ddlType string

	fs := flag.NewFlagSet("softdare", flag.ExitOnError)
	fs.StringVar(&ddlType, "type", "migrate", `
		Enter the type of the operation you want to perform.
		migrate: will create all the tables
		drop: will drop all the tables`)
	if err = fs.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parse flags error")
	}
	switch ddlType {
	case "migrate":
		if err = seed.MigrateTables(context.Background(), postgresDB); err != nil {
			return fmt.Errorf("migration error %s", err.Error())
		}
	case "drop":
		if err = seed.DropTables(context.Background(), postgresDB); err != nil {
			return fmt.Errorf("drop error %s", err.Error())
		}
	default:
		return errors.New("not valid flag type for action")
	}
	return nil
}
