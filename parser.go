package main

import (
	"awesomeProject1/repos"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Parser struct {
	EventsRepo  *repos.EventsRepository
	TablesRepo  *repos.TablesRepository
	ClientQueue *repos.ClientQueue
	ClubRepo    *repos.ClubPopulation
}

// NewParser создает новый экземпляр парсера с указанными репозиториями и очередью клиентов.
func NewParser(eventsRepo *repos.EventsRepository, tablesRepo *repos.TablesRepository, clientQueue *repos.ClientQueue, clubRepo *repos.ClubPopulation) *Parser {
	return &Parser{
		EventsRepo:  eventsRepo,
		TablesRepo:  tablesRepo,
		ClientQueue: clientQueue,
		ClubRepo:    clubRepo,
	}
}

// readInputFile читает содержимое файла и возвращает строки.
func readInputFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// parseInput парсит строки входных данных и возвращает соответствующие значения.
func parseInput(lines []string) (int, time.Time, time.Time, int, []repos.Event, error) {
	// Парсинг количества столов
	numTables, err := strconv.Atoi(lines[0])
	if err != nil || numTables <= 0 {
		log.Fatalf("Ошибка в строке 1: некорректное количество столов\n")
	}

	// Парсинг времени работы клуба
	workingHours := strings.Split(lines[1], " ")
	if len(workingHours) != 2 {
		log.Fatalf("Ошибка в строке 2: некорректное время работы клуба\n")
	}
	startTime, err := time.Parse("15:04", workingHours[0])
	if err != nil {
		log.Fatalf("Ошибка в строке 2: некорректное время начала работы\n")
	}
	endTime, err := time.Parse("15:04", workingHours[1])
	if err != nil {
		log.Fatalf("Ошибка в строке 2: некорректное время окончания работы\n")
	}
	// Парсинг стоимости часа
	hourlyRate, err := strconv.Atoi(lines[2])
	if err != nil || hourlyRate <= 0 {
		log.Fatalf("Ошибка в строке 3: некорректная стоимость часа\n")
	}

	// Парсинг событий
	eventStrings := lines[3:]

	var events []repos.Event

	for i, str := range eventStrings {
		parts := strings.Split(str, " ")
		if len(parts) < 3 {
			log.Fatalf("Ошибка в строке %d: некорректный формат события\n", i+4)
		}
		eventTime, err := time.Parse("15:04", parts[0])
		if err != nil {
			log.Fatalf("Ошибка в строке %d: некорректное время события\n", i+4)
		}
		if i > 0 && eventTime.Before(events[i-1].Time) {
			log.Fatalf("Ошибка в строке %d: события не идут последовательно во времени\n", i+4)
		}
		eventType := parts[1]
		clientName := parts[2]
		for _, r := range clientName {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
				log.Fatalf("Ошибка в строке %d: некорректное имя клиента\n", i+4)
			}
		}
		details := ""
		if len(parts) > 3 {
			details = strings.Join(parts[3:], " ")
		}

		event := repos.Event{
			Time:       eventTime,
			Type:       eventType,
			ClientName: clientName,
			Details:    details,
		}

		events = append(events, event)
	}
	return numTables, startTime, endTime, hourlyRate, events, nil
}
func (p *Parser) Run(fileName string) error {
	lines, err := readInputFile(fileName)
	if err != nil {
		return err
	}

	// Парсинг входных данных
	numTables, startTime, endTime, hourlyRate, events, err := parseInput(lines)
	if err != nil {
		return err
	}
	p.TablesRepo.AddTables(numTables)
	for _, event := range events {
		HandleEvent(event, p.EventsRepo, p.TablesRepo, p.ClientQueue, p.ClubRepo, startTime, endTime, hourlyRate)
	}
	handleEndOfDay(p.EventsRepo, p.TablesRepo, hourlyRate, endTime)
	formattedTime := startTime.Format("15:04")
	fmt.Printf("%s\n", formattedTime)
	//fmt.Println(startTime, endTime, hourlyRate, events)
	for _, event := range p.EventsRepo.Events {
		formattedTime := event.Time.Format("15:04")
		fmt.Printf("%s %s %s %s\n", formattedTime, event.Type, event.ClientName, event.Details)
	}
	formattedTime = endTime.Format("15:04")
	fmt.Printf("%s\n", formattedTime)
	for _, table := range p.TablesRepo.Tables {
		hours := int(table.TotalTimeOccupied.Hours())
		minutes := int(table.TotalTimeOccupied.Minutes()) % 60
		timeFormatted := fmt.Sprintf("%02d:%02d", hours, minutes)
		fmt.Printf("%d %d %s\n", table.ID, table.TotalAmountCollected, timeFormatted)
	}
	return nil
}
