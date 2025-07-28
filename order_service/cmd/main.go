package main

import (
	_ "github.com/lib/pq"
	"github.com/mayberries/shared/pkg/log"
	initApp "order_service/internal/init"
	"os"
)

func main() {
	err := initApp.App()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
