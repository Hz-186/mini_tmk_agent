package main

import (
	"log"
	"os"
	"project_for_tmk_04_06/internal/config"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
