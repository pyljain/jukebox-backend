package main

import (
	"flag"
	"jukebox/internal/db"
	"jukebox/internal/server"
	"log"

	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	port := flag.Int("p", 8090, "Port to listen on")
	dbConnectionString := flag.String("db", "jukebox.db", "Database connection string")
	driver := flag.String("driver", "sqlite3", "Database driver")
	flag.Parse()

	jubeboxDB, err := db.NewDB(*driver, *dbConnectionString)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer jubeboxDB.Close()

	llm, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	s := server.New(*port, jubeboxDB, llm)

	err = s.Start()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
