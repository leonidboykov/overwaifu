package overwaifu

// Score holds all scores
type Score struct {
	All                 int     `json:"all" bson:"all"`
	Safe                int     `json:"safe" bson:"safe"`
	Questionable        int     `json:"questionable" bson:"questionable"`
	Explicit            int     `json:"explicit" bson:"explicit"`
	GenderSwaps         int     `json:"gender_swaps" bson:"gender_swaps"`
	Lewd                float64 `json:"lewd" bson:"lewd"`
	Pure                float64 `json:"pure" bson:"pure"`
	VirginKillerSweater int     `json:"virgin_killer_sweater" bson:"virgin_killer_sweater"`
	Selfie              int     `json:"selfie" bson:"selfie"`
}
