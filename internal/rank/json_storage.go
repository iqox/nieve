package rank

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONStorage struct {
	path string
	mu   sync.Mutex
	data map[string]int
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	s := &JSONStorage{path: path, data: make(map[string]int)}
	if b, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(b, &s.data)
	}
	return s, nil
}

func (s *JSONStorage) PlayerRank(xuid string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.data[xuid]; ok {
		return id, nil
	}
	return -1, nil
}

func (s *JSONStorage) SavePlayerRank(xuid string, rankID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[xuid] = rankID
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0644)
}
