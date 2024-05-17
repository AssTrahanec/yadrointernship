package repos

type ClubPopulation struct {
	clientName []string
}

func NewClubPopulation() *ClubPopulation {
	return &ClubPopulation{}
}
func (cp *ClubPopulation) IsClientInClub(name string) bool {
	for _, client := range cp.clientName {
		if client == name {
			return true
		}
	}
	return false
}

func (cp *ClubPopulation) AddClient(clientName string) {
	cp.clientName = append(cp.clientName, clientName)
}

func (cp *ClubPopulation) RemoveClient(clientName string) {
	for i, name := range cp.clientName {
		if name == clientName {
			cp.clientName = append(cp.clientName[:i], cp.clientName[i+1:]...)
			return
		}
	}
}
