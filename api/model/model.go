package model

type ResponseJson struct {
	Responses Location `json:"response"`
}

type Location struct {
	Location []LocationDetail `json:"location"`
}

type LocationDetail struct {
	Postal     string `json:"postal"`
}

type ResponseErrorJson struct {
	Response ErrorDetail `json:"response"`
}

type ErrorDetail struct {
	Error string `json:"error"`
}