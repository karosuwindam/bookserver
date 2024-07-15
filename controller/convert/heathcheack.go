package convert

import (
	"sync"
	"time"
)

type dataStore struct {
	Status     bool
	Startfile  map[string]string
	Endfile    map[string]string
	StatusTIme map[string]time.Time
	mu         sync.Mutex
}

type DataStoreOutput struct {
	Status    bool     `json:"Status"`
	Startfile []string `json:"Startfile"`
	Endfile   []string `json:"Endfile"`
}

func CheackHealth() DataStoreOutput {
	statusData.mu.Lock()
	defer statusData.mu.Unlock()
	output := DataStoreOutput{
		Status:    statusData.Status,
		Startfile: []string{},
		Endfile:   []string{},
	}
	for _, r := range statusData.Startfile {
		output.Startfile = append(output.Startfile, r)
	}
	for _, r := range statusData.Endfile {
		output.Endfile = append(output.Endfile, r)
	}
	return output
}

var statusData dataStore = dataStore{}

func DataStoreInit() error {
	statusData.Startfile = make(map[string]string)
	statusData.Endfile = make(map[string]string)
	statusData.StatusTIme = make(map[string]time.Time)
	return nil
}

func (t *dataStore) On() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Status = true
}
func (t *dataStore) Off() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Status = false
}

func (t *dataStore) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for s, ts := range t.StatusTIme {
		if time.Now().Sub(ts) > 5*time.Minute { //5分を超えたとき
			delete(t.Startfile, s)
			delete(t.Endfile, s)
			delete(t.StatusTIme, s)
		}
	}
	// t.Startfile = make(map[string]string)
	// t.Endfile = make(map[string]string)
}

func (t *dataStore) Add(s string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Startfile[s] = s
	t.StatusTIme[s] = time.Now()
}

func (t *dataStore) Change(s string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.Startfile, s)
	t.Endfile[s] = s
	t.StatusTIme[s] = time.Now()
}

//状態を返すための処理プログラム
