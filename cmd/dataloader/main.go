package main

import (
	"log"

	"accounts/tools/dataloader"
)

func main() {
	if err := dataloader.Run(); err != nil {
		log.Fatal(err)
	}
}
