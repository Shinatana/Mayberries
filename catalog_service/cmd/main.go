package main

import (
	initApp "catalog_service/internal/init"
	"catalog_service/pkg/log"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	err := initApp.App()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
