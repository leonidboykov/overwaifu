package model

// Score ...
type Score struct {
	All          int `json:"all"`
	Safe         int `json:"safe"`
	Questionable int `json:"questionable"`
	Explicit     int `json:"explicit"`
}
