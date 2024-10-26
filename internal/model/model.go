package model

type Message struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	User               string  `json:"user"`
	State              string  `json:"state"`
	StartTime          string  `json:"startTime"`
	FinishedTime       string  `json:"finishedTime"`
	Duration           string  `json:"duration"`
	ErrorMessage       string  `json:"errorMessage"`
	PercentageComplete float64 `json:"percentageComplete"`
	Color              string  `json:"color"`
}
