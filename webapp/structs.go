package webapp

type Measurement struct {
	Id          string  `json:"id"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	Voltage     float64 `json:"voltage"`
	Timestamp   int     `json:"timestamp"`
}

type ChartAction struct {
	Measurements []Measurement
}
