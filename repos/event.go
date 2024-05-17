package repos

import "time"

// Event представляет событие в компьютерном клубе.
type Event struct {
	Time       time.Time
	Type       string
	ClientName string
	Details    string
}

// EventsRepository представляет репозиторий для управления событиями.
type EventsRepository struct {
	Events []Event
}

// NewEventsRepository создает новый репозиторий событий.
func NewEventsRepository() *EventsRepository {
	return &EventsRepository{}
}

// AddEvent добавляет событие в репозиторий.
func (repo *EventsRepository) AddEvent(event Event) {
	repo.Events = append(repo.Events, event)
}
