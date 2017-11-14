package overwaifu

// Score holds all scores
type Score struct {
	All                 int     `json:"all"`
	Safe                int     `json:"safe"`
	Questionable        int     `json:"questionable"`
	Explicit            int     `json:"explicit"`
	GenderSwaps         int     `json:"gender_swaps"`
	Lewd                float64 `json:"lewd"`
	Pure                float64 `json:"pure"`
	VirginKillerSweater int     `json:"virgin_killer_sweater"`
	Selfie              int     `json:"selfie"`
}
