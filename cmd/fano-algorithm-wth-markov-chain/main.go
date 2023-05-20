package main

import (
	"fano-algorithm/internal/pkg/app"
	"log"
)

func main() {
	err := app.WorkWithIO()
	if err != nil {
		log.Fatal(err)
	}
}
