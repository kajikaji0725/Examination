package model

import "time"

type ResponseJson struct {
	Response Location `json:"response"`
}

type Location struct {
	Error    string           `json:"error"`
	Location []LocationDetail `json:"location"`
}

type LocationDetail struct {
	Postal string    `json:"postal"`
	Date   time.Time `json:"date"`
}

type AccessLogs struct {
	AccessLog []AccessLog `json:"access_log"`
}

type AccessLog struct {
	Postal       string `json:"postal_code"`
	RequestCount int    `json:"RequestCount"`
}
