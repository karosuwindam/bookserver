package healthmessage

import (
	"encoding/json"
	"sync"
)

type healthMessageMu struct {
	Name    string `json:"Name"`
	Flag    bool   `json:"Flag"`
	Message string `json:"Message"`
	mu      sync.Mutex
}

type HealthMessage struct {
	Name    string `json:"Name"`
	Flag    bool   `json:"Flag"`
	Message string `json:"Message"`
}

func Create(name string) healthMessageMu {
	t := healthMessageMu{}
	t.Name = name
	return t
}

func (t *healthMessageMu) ChangeMessage(vs ...interface{}) {
	for _, v := range vs {
		switch v.(type) {
		case string:
			t.mu.Lock()
			t.Message = v.(string)
			t.mu.Unlock()
		case bool:
			t.mu.Lock()
			t.Flag = v.(bool)
			t.mu.Unlock()
		}
	}
}

func (t *healthMessageMu) JsonOut() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	tmp, _ := json.Marshal(t)
	return string(tmp)
}

func (t *healthMessageMu) ChangeOut() HealthMessage {
	t.mu.Lock()
	defer t.mu.Unlock()
	return HealthMessage{Name: t.Name, Flag: t.Flag, Message: t.Message}
}
