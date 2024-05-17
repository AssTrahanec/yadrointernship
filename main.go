package main

import (
	"awesomeProject1/repos"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./task.exe <filename>")
	}

	filename := os.Args[1]

	eventsRepository := repos.NewEventsRepository()
	clientQueue := repos.NewClientQueue()
	tablesRepository := repos.NewTablesRepository()
	clubRepository := repos.NewClubPopulation()
	parser := NewParser(eventsRepository, tablesRepository, clientQueue, clubRepository)
	err := parser.Run(filename)
	if err != nil {
		return
	}

}
