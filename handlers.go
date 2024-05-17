package main

import (
	"awesomeProject1/repos"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

func HandleEvent(event repos.Event, eventsRepo *repos.EventsRepository, tablesRepo *repos.TablesRepository, clientQueue *repos.ClientQueue,
	clubRepo *repos.ClubPopulation, startTime time.Time, endTime time.Time, hourlyRate int) {
	switch event.Type {
	case "1":
		HandleClientArrived(event, eventsRepo, tablesRepo, clientQueue, clubRepo, startTime, endTime)
	case "2":
		HandleClientReserveTable(event, eventsRepo, tablesRepo, clientQueue, clubRepo, hourlyRate)
	case "3":
		HandleClientWaiting(event, eventsRepo, tablesRepo, clientQueue)
	case "4":
		HandleClientLeft(event, eventsRepo, tablesRepo, clientQueue, clubRepo, hourlyRate)
	default:
		//return
	}
}

func HandleClientArrived(event repos.Event, eventsRepo *repos.EventsRepository, tablesRepo *repos.TablesRepository, clientQueue *repos.ClientQueue,
	clubRepo *repos.ClubPopulation, startTime time.Time, endTime time.Time) {
	eventsRepo.AddEvent(event)
	// Проверка времени работы клуба
	clubRepo.AddClient(event.ClientName)
	if !(event.Time.After(startTime) && event.Time.Before(endTime)) {
		eventsRepo.AddEvent(repos.Event{Time: event.Time, Type: "13", ClientName: event.ClientName, Details: "NotOpenYet"})
		return
	}
	// Проверка есть ли клиент в клубе
	if clientQueue.IsClientInQueue(event.ClientName) || tablesRepo.IsClientInside(event.ClientName) {
		eventsRepo.AddEvent(repos.Event{Time: event.Time, Type: "13", ClientName: event.ClientName, Details: "YouShallNotPass"})
		return
	}
	return
}
func HandleClientReserveTable(event repos.Event, eventsRepo *repos.EventsRepository, tablesRepo *repos.TablesRepository, clientQueue *repos.ClientQueue,
	clubRepo *repos.ClubPopulation, hourlyRate int) {
	eventsRepo.AddEvent(event)
	// Проверка есть ли клиент в клубе
	//if !(clientQueue.IsClientInQueue(event.ClientName) || tablesRepo.IsClientInside(event.ClientName)) {
	if !clubRepo.IsClientInClub(event.ClientName) {
		eventsRepo.AddEvent(repos.Event{Time: event.Time, Type: "13", ClientName: event.ClientName, Details: "ClientUnknown"})
		return
	}

	tableNumber, err := strconv.Atoi(event.Details)
	if err != nil {
		fmt.Println(event)
		os.Exit(1)
	}

	// Проверка свободен ли выбранный стол
	if !tablesRepo.IsTableFree(tableNumber) {
		eventsRepo.AddEvent(repos.Event{Time: event.Time, Type: "13", ClientName: event.ClientName, Details: "PlaceIsBusy"})
		return
	}
	// Освободить стол, если меняет место
	if tablesRepo.IsClientInside(event.ClientName) {
		tablesRepo.FreeTable(event.ClientName, event.Time, hourlyRate)
	}

	// Занял стол
	tablesRepo.ReserveTable(tableNumber, event.ClientName, event.Time, hourlyRate)
	return
}
func HandleClientWaiting(event repos.Event, eventsRepo *repos.EventsRepository, tablesRepo *repos.TablesRepository, clientQueue *repos.ClientQueue) {
	eventsRepo.AddEvent(event)
	if tablesRepo.CheckFreeSpace() != 0 {
		eventsRepo.AddEvent(repos.Event{Time: event.Time, Type: "13", ClientName: event.ClientName, Details: "ICanWaitNoLonger"})
		return
	}
	if tablesRepo.GetTablesNum() < clientQueue.QueueLength() {
		eventsRepo.AddEvent(repos.Event{Time: event.Time, Type: "11", ClientName: event.ClientName, Details: ""})
		return
	}
	clientQueue.AddClient(event.ClientName)
	return
}
func HandleClientLeft(event repos.Event, eventsRepo *repos.EventsRepository, tablesRepo *repos.TablesRepository, clientQueue *repos.ClientQueue,
	clubRepo *repos.ClubPopulation, hourlyRate int) {
	eventsRepo.AddEvent(event)
	// Проверка есть ли клиент в клубе
	if !(clientQueue.IsClientInQueue(event.ClientName) || tablesRepo.IsClientInside(event.ClientName)) {
		eventsRepo.AddEvent(repos.Event{Time: event.Time, Type: "13", ClientName: event.ClientName, Details: "ClientUnknown"})
		return
	}
	// Освободил стол, если занимал
	if tablesRepo.IsClientInside(event.ClientName) {
		freeTableId := tablesRepo.FreeTable(event.ClientName, event.Time, hourlyRate)

		// Проверка очереди
		if clientQueue.QueueLength() > 0 {
			newClient := clientQueue.GetNextClient()
			tablesRepo.ReserveTable(freeTableId, newClient, event.Time, hourlyRate)
			eventsRepo.AddEvent(repos.Event{Time: event.Time, Type: "12", ClientName: event.ClientName, Details: strconv.Itoa(freeTableId)})
			return
		}
	}
	clubRepo.RemoveClient(event.ClientName)
	return
}
func handleEndOfDay(eventsRepo *repos.EventsRepository, tablesRepo *repos.TablesRepository, hourlyRate int, endTime time.Time) {
	if tablesRepo.GetClientsOnTables() != nil {
		clients := tablesRepo.GetClientsOnTables()
		// Сортировка по алфавиту
		sort.Strings(clients)
		for _, client := range clients {
			tablesRepo.FreeTable(client, endTime, hourlyRate)
			eventsRepo.AddEvent(repos.Event{Time: endTime, Type: "11", ClientName: client, Details: ""})
		}
	}
}
