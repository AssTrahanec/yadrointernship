package repos

import (
	"math"
	"time"
)

type Table struct {
	ID                   int
	IsFree               bool
	clientName           string
	StartTime            time.Time
	TotalTimeOccupied    time.Duration
	TotalAmountCollected int // Общая сумма денег, собранная за использование стола
}

type TablesRepository struct {
	Tables []*Table
}

func NewTablesRepository() *TablesRepository {
	return &TablesRepository{}
}
func (repo *TablesRepository) AddTables(numTables int) {
	tables := make([]*Table, numTables)
	for i := 0; i < numTables; i++ {

		tables[i] = &Table{ID: i + 1, IsFree: true, clientName: "", TotalTimeOccupied: 0, TotalAmountCollected: 0}
	}
	repo.Tables = tables
}
func (repo *TablesRepository) CheckFreeSpace() int {
	for _, table := range repo.Tables {
		if table.IsFree == true {
			return table.ID
		}
	}
	return 0
}
func (repo *TablesRepository) IsTableFree(id int) bool {
	for _, table := range repo.Tables {
		if table.ID == id && table.IsFree {
			return true
		}
	}
	return false
}

// ReserveTable занять стол
func (repo *TablesRepository) ReserveTable(id int, clientName string, startTime time.Time, hourlyRate int) {
	for _, table := range repo.Tables {
		if table.ID == id {
			table.IsFree = false
			table.clientName = clientName
			table.StartTime = startTime

		}
	}
}

// FreeTable освободить стол, возвращает id освобожденного стола
func (repo *TablesRepository) FreeTable(clientName string, endTime time.Time, hourlyRate int) int {
	for _, table := range repo.Tables {
		if table.clientName == clientName {
			duration := endTime.Sub(table.StartTime)
			hours := duration.Hours()
			roundedHours := math.Ceil(hours)

			amountCollected := int(roundedHours) * hourlyRate
			table.TotalTimeOccupied += duration
			table.TotalAmountCollected += amountCollected
			table.IsFree = true
			table.clientName = ""
			table.StartTime = time.Time{}
			return table.ID
		}
	}
	return 0
}

// IsClientInside возвращает сидит ли клиент за столом
func (repo *TablesRepository) IsClientInside(clientName string) bool {
	for _, table := range repo.Tables {
		if table.clientName == clientName {
			return true
		}
	}
	return false
}

// GetTablesNum возвращает количество столов
func (repo *TablesRepository) GetTablesNum() int {
	return len(repo.Tables)
}

// GetClientsOnTables возвращает имена клиентов сидящих за столами
func (repo *TablesRepository) GetClientsOnTables() []string {
	var clients []string
	for _, table := range repo.Tables {
		if !table.IsFree {
			clients = append(clients, table.clientName)
		}
	}
	return clients
}
