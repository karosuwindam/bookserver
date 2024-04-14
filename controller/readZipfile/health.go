package readzipfile

//動作状態の確認

type HealthData struct {
	Status   bool     `json:"Status"`
	Name     []string `json:"Name"`
	CashSize int64    `json:"Cashsize"`
}

func HealthCheck() HealthData {
	output := HealthData{
		Status:   loopflag,
		Name:     []string{},
		CashSize: 0,
	}
	dataStore.mu.Lock()
	tmps := dataStore.cashZipSize
	dataStore.mu.Unlock()
	for s, i := range tmps {
		output.CashSize += int64(i)
		output.Name = append(output.Name, s)
	}
	return output
}
