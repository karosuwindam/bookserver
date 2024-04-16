package copyfile

type HealthCopyFIleData struct {
	Status bool `json:"Status"`
}

func HealthCheck() HealthCopyFIleData {
	output := HealthCopyFIleData{
		Status: loopflag,
	}
	return output
}
