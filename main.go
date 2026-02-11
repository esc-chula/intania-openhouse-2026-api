package main

import (
	"log"

	"github.com/esc-chula/intania-openhouse-2026-api/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
