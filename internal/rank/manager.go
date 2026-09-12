package rank

import (
	"fmt"
	"sync"
)

type Manager struct {
	mu       sync.RWMutex
	ranks    map[int]*Rank
	players  map[string]int
	defaultR int
	store    Storage
}

func NewManager(defaultRankID int, store Storage) *Manager {
	return &Manager{
		ranks:    make(map[int]*Rank),
		players:  make(map[string]int),
		defaultR: defaultRankID,
		store:    store,
	}
}

func (m *Manager) RegisterRank(r *Rank) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ranks[r.ID] = r
}

func (m *Manager) RankByID(id int) *Rank {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if r, ok := m.ranks[id]; ok {
		return r
	}
	return m.ranks[m.defaultR]
}

func (m *Manager) PlayerRank(xuid string) *Rank {
	m.mu.RLock()
	id, ok := m.players[xuid]
	m.mu.RUnlock()
	if !ok {
		return m.ranks[m.defaultR]
	}
	return m.RankByID(id)
}

func (m *Manager) SetPlayerRank(xuid string, rankID int) error {
	m.mu.Lock()
	if _, ok := m.ranks[rankID]; !ok {
		m.mu.Unlock()
		return fmt.Errorf("rank %q doesn't exist", rankID)
	}
	m.players[xuid] = rankID
	m.mu.Unlock()

	return m.store.SavePlayerRank(xuid, rankID)
}

func (m *Manager) LoadPlayer(xuid string) error {
	id, err := m.store.PlayerRank(xuid)
	if err != nil {
		return err
	}
	m.mu.Lock()
	if id == -1 {
		id = m.defaultR
	}
	m.players[xuid] = id
	m.mu.Unlock()
	return nil
}
