package repos

type ClientQueue struct {
	clientName []string
}

// NewClientQueue создает новый экземпляр репозитория очереди клиентов.
func NewClientQueue() *ClientQueue {
	return &ClientQueue{}
}

// AddClient добавляет имя клиента в очередь.
func (cq *ClientQueue) AddClient(clientName string) {
	cq.clientName = append(cq.clientName, clientName)
}

// GetNextClient извлекает и возвращает первого клиента из очереди, если она не пуста.
func (cq *ClientQueue) GetNextClient() string {
	if len(cq.clientName) == 0 {
		return "" // Если очередь пуста, возвращаем пустую строку
	}
	// Извлекаем первого клиента из очереди
	clientName := cq.clientName[0]
	cq.clientName = cq.clientName[1:]
	return clientName
}

// IsClientInQueue возвращает находится ли клиент в очереди
func (cq *ClientQueue) IsClientInQueue(clientName string) bool {
	for _, client := range cq.clientName {
		if client == clientName {
			return true
		}
	}
	return false
}

// QueueLength возвращает длину очереди
func (cq *ClientQueue) QueueLength() int {
	return len(cq.clientName)
}
