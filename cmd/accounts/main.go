package main

import (
	"log"

	"accounts/app"
)

func main() {
	log.Fatal(app.Serve())
}
