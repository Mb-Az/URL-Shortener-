package mapstruct

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"urlShortener/mapInterface"
)

type MapS struct {
	mapS map[string]string
	mu   sync.RWMutex
}

type MapSync struct {
	mapS sync.Map
	mu   sync.Mutex
}

func (m *MapS) Get(short string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.mapS[short]
	if !ok {
		return "", fmt.Errorf("value not found")
	}
	return value, nil
}

func (m *MapS) Add(short, url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.mapS[short]
	if ok {
		return fmt.Errorf("This short-url is taken")
	}
	m.mapS[short] = url
	return nil
}

func (m *MapS) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	jsonString, err := json.MarshalIndent(m.mapS, "", " ")
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	file, err := os.OpenFile("map.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("Error opening the file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonString)
	if err != nil {
		return fmt.Errorf("Error writing file: %w", err)
	}
	// fmt.Println("File saved")
	return nil

}

func (m *MapS) Load() {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := os.ReadFile("map.json")
	if err != nil {
		return
	}

	if err := json.Unmarshal(data, &m.mapS); err != nil {
		return
	}

}

func (m *MapSync) Get(short string) (string, error) {
	// m.RLock()
	val, ok := m.mapS.Load(short)
	if !ok {
		return "", fmt.Errorf("Not found")
	}
	if val, ok := val.(string); ok {
		return val, nil
	}
	return "", fmt.Errorf("Could not convert to string")

}

func (m *MapSync) Add(short, url string) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	// _, ok := m.mapS[short]
	_, ok := m.mapS.Load(short)
	if ok {
		return fmt.Errorf("This short-url is taken")
	}
	m.mapS.Store(short, url)

	return nil
}

func (m *MapSync) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mapNormal := make(map[string]string, 50)
	m.mapS.Range(func(key, value any) bool {

		if key, ok := key.(string); ok {
			if value, ok := value.(string); ok {
				mapNormal[key] = value
				return true
			}
			return false
		}

		return false
	})
	jsonString, err := json.MarshalIndent(mapNormal, "", " ")
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	file, err := os.OpenFile("map.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("Error opening the file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonString)
	if err != nil {
		return fmt.Errorf("Error writing file: %w", err)
	}
	// fmt.Println("File saved")
	return nil

}

func (m *MapSync) Load() {
	data, err := os.ReadFile("map.json")
	mapNormal := make(map[string]string, 50)

	if err != nil {
		return
	}

	if err := json.Unmarshal(data, &mapNormal); err != nil {
		return
	}
	for key, value := range mapNormal {
		m.mapS.Store(key, value)
	}

}

func New(option string) mapInterface.Map {
	switch option {

	// case "RWlock":
	// 	return &MapS{mapS: make(map[string]string, 100)}

	case "sync":
		return &MapSync{mapS: sync.Map{}}

	default:
		return &MapS{mapS: make(map[string]string, 100)}
	}

}
