package main

import (
	initApp "auth_service/internal/init"
	"auth_service/pkg/log"
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
