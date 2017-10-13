package overwaifu

// Score holds all scores
type Score struct {
	All          int `json:"all"`
	Safe         int `json:"safe"`
	Questionable int `json:"questionable"`
	Explicit     int `json:"explicit"`
	Genderswaps  int `json:"gender_swaps"`
}
