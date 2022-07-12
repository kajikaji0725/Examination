package model

type ResponseJson struct {
	Response Location `json:"response"`
}

type Location struct {
	Error    string           `json:"error"`
	Location []LocationDetail `json:"location"`
}

type LocationDetail struct {
	City       string `json:"city"`
	Prefecture string `json:"prefecture"`
	Postal     string `json:"postal"`
}
